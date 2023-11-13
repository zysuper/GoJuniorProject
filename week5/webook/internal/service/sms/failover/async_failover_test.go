package failover

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	repomocks "gitee.com/geekbang/basic-go/webook/internal/repository/mocks"
	cb "gitee.com/geekbang/basic-go/webook/internal/service/sms/circuit_breaker"
	smsmocks "gitee.com/geekbang/basic-go/webook/internal/service/sms/mocks"
	"gitee.com/geekbang/basic-go/webook/pkg/limiter"
	limitermocks "gitee.com/geekbang/basic-go/webook/pkg/limiter/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestAsyncFailoverService_Send(t *testing.T) {
	tests := []struct {
		name      string
		retryTime time.Duration
		retryCnt  int
		wantErr   error
		mock      func(controller *gomock.Controller) (cb.CircuitBreaker, limiter.Limiter, repository.MsgRepository)
	}{
		{
			name:      "直接成功",
			retryTime: time.Millisecond * 30,
			retryCnt:  3,
			mock: func(controller *gomock.Controller) (cb.CircuitBreaker, limiter.Limiter, repository.MsgRepository) {
				lt := limitermocks.NewMockLimiter(controller)
				sms := smsmocks.NewMockService(controller)
				cbb := cb.NewCircuitBreaker(3, time.Second, cb.NewSmsCircuitBreakerAdapter(sms))
				msg := repomocks.NewMockMsgRepository(controller)
				sms.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				lt.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(false, nil)
				return cbb, lt, msg
			},
		},
		{
			name:      "被熔断, 然后过一会，熔断器进入 half 状态，再发送成功了",
			retryTime: time.Millisecond * 30,
			retryCnt:  3,
			mock: func(controller *gomock.Controller) (cb.CircuitBreaker, limiter.Limiter, repository.MsgRepository) {
				lt := limitermocks.NewMockLimiter(controller)
				cbb := smsmocks.NewMockCircuitBreaker(controller)
				msg := repomocks.NewMockMsgRepository(controller)
				// sms.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				lt.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(false, nil)
				cbb.EXPECT().Do(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(cb.CbCloseError)
				msg.EXPECT().Create(gomock.Any(), gomock.Any()).Return(int64(123), nil)
				msg.EXPECT().FindById(gomock.Any(), int64(123)).Return(domain.Msg{
					Id:      "123",
					TplId:   "123",
					Args:    []string{"hello,world"},
					Numbers: []string{"123q444"},
				}, nil)
				// 熔断后恢复了。
				cbb.EXPECT().Do(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return cbb, lt, msg
			},
			wantErr: cb.CbCloseError,
		},
		{
			name:      "被限流后，异步存储后，再重试发送成功了亲",
			retryTime: time.Millisecond * 30,
			retryCnt:  3,
			mock: func(controller *gomock.Controller) (cb.CircuitBreaker, limiter.Limiter, repository.MsgRepository) {
				lt := limitermocks.NewMockLimiter(controller)
				sms := smsmocks.NewMockService(controller)
				c := cb.NewCircuitBreaker(3, time.Second, cb.NewSmsCircuitBreakerAdapter(sms))
				msgRepo := repomocks.NewMockMsgRepository(controller)
				lt.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(true, nil)
				msgRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(int64(123), nil)
				msgRepo.EXPECT().FindById(gomock.Any(), int64(123)).Return(domain.Msg{
					Id:      "123",
					TplId:   "123",
					Args:    []string{"hello,world"},
					Numbers: []string{"123q444"},
				}, nil)
				sms.EXPECT().Send(gomock.Any(), "123", []string{"hello,world"}, gomock.Any()).Return(nil)
				return c, lt, msgRepo
			},
			wantErr: LimitedError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c, lt, msgRepo := tt.mock(ctrl)
			service := NewAsyncFailoverService(c, lt, msgRepo, tt.retryTime, tt.retryCnt)
			err := service.Send(context.Background(), "123", []string{"hello,world"}, "123q444")
			// 睡一会，让 go routine 充分燃烧.
			time.Sleep(time.Second)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
