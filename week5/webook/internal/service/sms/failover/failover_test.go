package failover

import (
	"context"
	"errors"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms"
	smsmocks "gitee.com/geekbang/basic-go/webook/internal/service/sms/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestFailOverSmsService_Send(t *testing.T) {
	tests := []struct {
		name    string
		mock    func(controller *gomock.Controller) []sms.SmsService
		wantErr error
	}{
		{
			name: "第一次发送成功",
			mock: func(controller *gomock.Controller) []sms.SmsService {
				svc := smsmocks.NewMockSmsService(controller)
				svc.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return []sms.SmsService{svc}
			},
			wantErr: nil,
		},
		{
			name: "第二次发送成功",
			mock: func(controller *gomock.Controller) []sms.SmsService {
				svc0 := smsmocks.NewMockSmsService(controller)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("居然发送失败了，呜呜"))
				svc1 := smsmocks.NewMockSmsService(controller)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return []sms.SmsService{svc0, svc1}
			},
			wantErr: nil,
		},
		{
			name: "全部发送失败",
			mock: func(controller *gomock.Controller) []sms.SmsService {
				svc0 := smsmocks.NewMockSmsService(controller)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("居然发送失败了，呜呜"))
				svc1 := smsmocks.NewMockSmsService(controller)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("又发送失败了"))
				return []sms.SmsService{svc0, svc1}
			},
			wantErr: AllSendFailed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewFailOverSmsService(tt.mock(ctrl))
			err := svc.Send(context.Background(), "123", []string{"123"}, "12233")
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
