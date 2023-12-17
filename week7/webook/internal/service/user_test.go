package service

import (
	"context"
	"errors"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	repomocks "gitee.com/geekbang/basic-go/webook/internal/repository/mocks"
	svcmocks "gitee.com/geekbang/basic-go/webook/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestPasswordEncrypt(t *testing.T) {
	password := []byte("123456#hello")
	encrypted, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	assert.NoError(t, err)
	println(string(encrypted))
	err = bcrypt.CompareHashAndPassword(encrypted, []byte("123456#hello"))
	assert.NoError(t, err)
}

func Test_userService_Login(t *testing.T) {
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) (repository.UserRepository, PasswordValidateService)

		// 预期输入
		ctx      context.Context
		email    string
		password string

		wantUser domain.User
		wantErr  error
	}{
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) (repository.UserRepository, PasswordValidateService) {
				repo := repomocks.NewMockUserRepository(ctrl)
				pv := svcmocks.NewMockPasswordValidateService(ctrl)
				pv.EXPECT().ComparePassword("$2a$10$.l0JHmM7a2PdJ.A9gsmVyerEDlp1WhxsglC34S4UJH4TuHhWY7Tfq", gomock.Any())
				repo.EXPECT().
					FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{
						Email: "123@qq.com",
						// 你在这边拿到的密码，就应该是一个正确的密码
						// 加密后的正确的密码
						Password: "$2a$10$.l0JHmM7a2PdJ.A9gsmVyerEDlp1WhxsglC34S4UJH4TuHhWY7Tfq",
						Phone:    "15212345678",
					}, nil)
				return repo, pv
			},
			email: "123@qq.com",
			// 用户输入的，没有加密的
			password: "123456#hello",

			wantUser: domain.User{
				Email:    "123@qq.com",
				Password: "$2a$10$.l0JHmM7a2PdJ.A9gsmVyerEDlp1WhxsglC34S4UJH4TuHhWY7Tfq",
				Phone:    "15212345678",
			},
		},

		{
			name: "用户未找到",
			mock: func(ctrl *gomock.Controller) (repository.UserRepository, PasswordValidateService) {
				repo := repomocks.NewMockUserRepository(ctrl)
				pv := svcmocks.NewMockPasswordValidateService(ctrl)
				repo.EXPECT().
					FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{}, repository.ErrUserNotFound)
				return repo, pv
			},
			email: "123@qq.com",
			// 用户输入的，没有加密的
			password: "123456#hello",
			wantErr:  ErrInvalidUserOrPassword,
		},

		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) (repository.UserRepository, PasswordValidateService) {
				repo := repomocks.NewMockUserRepository(ctrl)
				pv := svcmocks.NewMockPasswordValidateService(ctrl)
				repo.EXPECT().
					FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{}, errors.New("db错误"))
				return repo, pv
			},
			email: "123@qq.com",
			// 用户输入的，没有加密的
			password: "123456#hello",
			wantErr:  errors.New("db错误"),
		},

		{
			name: "密码不对",
			mock: func(ctrl *gomock.Controller) (repository.UserRepository, PasswordValidateService) {
				repo := repomocks.NewMockUserRepository(ctrl)
				pv := svcmocks.NewMockPasswordValidateService(ctrl)
				pv.EXPECT().ComparePassword("$2a$10$.l0JHmM7a2PdJ.A9gsmVyerEDlp1WhxsglC34S4UJH4TuHhWY7Tfq", gomock.Any()).Return(errors.New("用户不存在或者密码不对"))
				repo.EXPECT().
					FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{
						Email: "123@qq.com",
						// 你在这边拿到的密码，就应该是一个正确的密码
						// 加密后的正确的密码
						Password: "$2a$10$.l0JHmM7a2PdJ.A9gsmVyerEDlp1WhxsglC34S4UJH4TuHhWY7Tfq",
						Phone:    "15212345678",
					}, nil)
				return repo, pv
			},
			email: "123@qq.com",
			// 用户输入的，没有加密的
			password: "123456#helloABCde",

			wantErr: ErrInvalidUserOrPassword,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo, pv := tc.mock(ctrl)
			svc := NewUserService(repo, pv)
			user, err := svc.Login(tc.ctx, tc.email, tc.password)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, user)
		})
	}
}
