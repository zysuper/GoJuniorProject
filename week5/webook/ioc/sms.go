package ioc

import (
	"gitee.com/geekbang/basic-go/webook/internal/service/sms"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms/local"
)

func InitSms() sms.SmsService {
	return local.NewService()
}
