package service

import (
	"context"
	"fmt"
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms"
	"math/rand"
)

const (
	codeTplId = "1877556"
)

var ErrCodeSendTooMany = repository.ErrCodeSendTooMany

type CodeService interface {
	Send(ctx context.Context, biz, phone string) error
	Verify(ctx context.Context,
		biz, phone, inputCode string) (bool, error)
}

type codeService struct {
	repo repository.CodeRepository
	sms  sms.SmsService
}

func NewCodeService(repo repository.CodeRepository, sms sms.SmsService) CodeService {
	return &codeService{repo: repo, sms: sms}
}

func (c *codeService) Send(ctx context.Context, biz, phone string) error {
	code := codeGen()
	err := c.repo.Set(ctx, biz, phone, code)
	// 你在这儿，是不是要开始发送验证码了？
	if err != nil {
		return err
	}
	return c.sms.Send(ctx, codeTplId, []string{code}, phone)
}

func (c *codeService) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	ok, err := c.repo.Verify(ctx, biz, phone, inputCode)
	if err == repository.ErrCodeVerifyTooMany {
		// 相当于，我们对外面屏蔽了验证次数过多的错误，我们就是告诉调用者，你这个不对
		return false, nil
	}
	return ok, err
}

func codeGen() string {
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}
