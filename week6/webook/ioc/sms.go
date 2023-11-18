package ioc

import (
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms"
	cb "gitee.com/geekbang/basic-go/webook/internal/service/sms/circuit_breaker"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms/failover"
	"github.com/redis/go-redis/v9"
	"time"
)

func InitSms(cb cb.CircuitBreaker, cmdable redis.Cmdable, repo repository.MsgRepository) sms.Service {
	// return ratelimit.NewRateLimitSMSService(local.NewService(), limiter.NewRedisSlidingWindowLimiter(cmd, time.Second, 1000))
	return failover.NewAsyncFailoverService(cb, NewRedisSmsLimiter(cmdable), repo, time.Millisecond*30, 3)
}
