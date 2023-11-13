package ioc

import (
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms"
	cb "gitee.com/geekbang/basic-go/webook/internal/service/sms/circuit_breaker"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms/failover"
	"gitee.com/geekbang/basic-go/webook/pkg/limiter"
	"time"
)

func InitSms(cb cb.CircuitBreaker, limiter limiter.Limiter, repo repository.MsgRepository) sms.Service {
	// return ratelimit.NewRateLimitSMSService(local.NewService(), limiter.NewRedisSlidingWindowLimiter(cmd, time.Second, 1000))
	return failover.NewAsyncFailoverService(cb, limiter, repo, time.Millisecond*30, 3)
}
