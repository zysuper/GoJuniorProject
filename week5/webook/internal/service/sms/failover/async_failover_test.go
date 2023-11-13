package failover

import (
	"context"
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
				limiter := limitermocks.NewMockLimiter(controller)
				sms := smsmocks.NewMockService(controller)
				cb := cb.NewCircuitBreaker(3, time.Second, func(args ...any) error {
					return sms.Send(args[0].(context.Context), args[1].(string), args[2].([]string), args[3].([]string)...)
				})
				msg := repomocks.NewMockMsgRepository(controller)
				sms.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				limiter.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(false, nil)
				return cb, limiter, msg
			},
		},
		{
			name:      "被熔断",
			retryTime: time.Millisecond * 30,
			retryCnt:  3,
			mock: func(controller *gomock.Controller) (cb.CircuitBreaker, limiter.Limiter, repository.MsgRepository) {
				limiter := limitermocks.NewMockLimiter(controller)
				cbb := smsmocks.NewMockCircuitBreaker(controller)
				msg := repomocks.NewMockMsgRepository(controller)
				//sms.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				limiter.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(false, nil)
				cbb.EXPECT().Do(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(cb.CbCloseError)
				return cbb, limiter, msg
			},
			wantErr: cb.CbCloseError,
		},
		{
			name:      "被限流",
			retryTime: time.Millisecond * 30,
			retryCnt:  3,
			mock: func(controller *gomock.Controller) (cb.CircuitBreaker, limiter.Limiter, repository.MsgRepository) {
				limiter := limitermocks.NewMockLimiter(controller)
				sms := smsmocks.NewMockService(controller)
				cb := cb.NewCircuitBreaker(3, time.Second, func(args ...any) error {
					return sms.Send(args[0].(context.Context), args[1].(string), args[2].([]string), args[3].([]string)...)
				})
				msgRepo := repomocks.NewMockMsgRepository(controller)
				//sms.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				limiter.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(true, nil)
				//msgRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(int64(123), nil)
				return cb, limiter, msgRepo
			},
			wantErr: LimitedError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			cb, limiter, msgRepo := tt.mock(ctrl)
			service := NewAsyncFailoverService(cb, limiter, msgRepo, tt.retryTime, tt.retryCnt)
			err := service.Send(context.Background(), "123", []string{"hello,world"}, "123q444")
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
