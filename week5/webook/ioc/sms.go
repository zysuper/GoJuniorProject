package ioc

import (
	"gitee.com/geekbang/basic-go/webook/internal/service/sms"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms/local"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms/ratelimit"
	"gitee.com/geekbang/basic-go/webook/pkg/limiter"
	"github.com/redis/go-redis/v9"
	"time"
)

func InitSms(cmd redis.Cmdable) sms.Service {
	return ratelimit.NewRateLimitSMSService(local.NewService(), limiter.NewRedisSlidingWindowLimiter(cmd, time.Second, 1000))
}
