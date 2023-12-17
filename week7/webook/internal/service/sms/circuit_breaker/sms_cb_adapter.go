package cb

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms"
)

type SmsCircuitBreakerAdapter struct {
	sms sms.Service
}

func NewSmsCircuitBreakerAdapter(sms sms.Service) CircuitBreaker {
	return &SmsCircuitBreakerAdapter{sms: sms}
}

func (s *SmsCircuitBreakerAdapter) Do(args ...any) error {
	return s.sms.Send(args[0].(context.Context), args[1].(string), args[2].([]string), args[3].([]string)...)
}
