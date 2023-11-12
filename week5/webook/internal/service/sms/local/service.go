package local

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms"
	"log"
)

type Service struct {
}

func (s Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	log.Println("验证码是", args)
	return nil
}

func NewService() sms.SmsService {
	return &Service{}
}
