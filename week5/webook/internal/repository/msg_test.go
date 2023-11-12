package repository

import (
	"context"
	"errors"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
	daomocks "gitee.com/geekbang/basic-go/webook/internal/repository/dao/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func Test_msgRepository_Create(t *testing.T) {
	tests := []struct {
		name    string
		mock    func(controller *gomock.Controller) dao.MsgDAO
		msg     domain.Msg
		wantErr error
	}{
		{
			name: "创建成功",
			msg: domain.Msg{
				Id:      "123",
				TplId:   "aaaa",
				Args:    []string{"xxxx"},
				Numbers: []string{"122344"},
			},
			mock: func(controller *gomock.Controller) dao.MsgDAO {
				dao := daomocks.NewMockMsgDAO(controller)
				dao.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				return dao
			},
		},
		{
			name: "创建失败",
			msg: domain.Msg{
				Id:      "123",
				TplId:   "aaaa",
				Args:    []string{"xxxx"},
				Numbers: []string{"122344"},
			},
			mock: func(controller *gomock.Controller) dao.MsgDAO {
				d := daomocks.NewMockMsgDAO(controller)
				d.EXPECT().Create(gomock.Any(),
					dao.Msg{
						Args: `{"id":"123","tplId":"aaaa","args":["xxxx"],"numbers":["122344"]}`}).
					Return(errors.New("创建失败"))
				return d
			},
			wantErr: errors.New("创建失败"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			dao := tt.mock(ctrl)
			svc := NewMsgRepository(dao)
			err := svc.Create(context.Background(), tt.msg)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func Test_msgRepository_FindById(t *testing.T) {
	tests := []struct {
		name    string
		mock    func(controller *gomock.Controller) dao.MsgDAO
		uid     int64
		wantErr error
		wantMsg domain.Msg
	}{
		{
			name: "查询成功",
			wantMsg: domain.Msg{
				Id:      "123",
				TplId:   "aaaa",
				Args:    []string{"xxxx"},
				Numbers: []string{"122344"},
			},
			uid: 123,
			mock: func(controller *gomock.Controller) dao.MsgDAO {
				d := daomocks.NewMockMsgDAO(controller)
				d.EXPECT().FindById(gomock.Any(), int64(123)).Return(dao.Msg{
					Id:   int64(123),
					Args: `{"id": "123", "tplId": "aaaa", "args": ["xxxx"], "numbers": ["122344"]}`}, nil)
				return d
			},
		},
		{
			name: "查询失败",
			mock: func(controller *gomock.Controller) dao.MsgDAO {
				dd := daomocks.NewMockMsgDAO(controller)
				dd.EXPECT().FindById(gomock.Any(), gomock.Any()).Return(dao.Msg{}, errors.New("查询失败"))
				return dd
			},
			wantErr: errors.New("查询失败"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			dao := tt.mock(ctrl)
			svc := NewMsgRepository(dao)
			msg, err := svc.FindById(context.Background(), tt.uid)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantMsg, msg)
		})
	}
}
