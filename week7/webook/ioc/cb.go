package ioc

import (
	cb "gitee.com/geekbang/basic-go/webook/internal/service/sms/circuit_breaker"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms/local"
	"time"
)

func NewSmsCircuitBreaker() cb.CircuitBreaker {
	return cb.NewCircuitBreaker(3, time.Second, cb.NewSmsCircuitBreakerAdapter(local.NewService()))
}
