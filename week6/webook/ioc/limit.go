package ioc

import (
	"gitee.com/geekbang/basic-go/webook/pkg/limiter"
	"github.com/redis/go-redis/v9"
	"time"
)

func NewRedisSmsLimiter(cmdable redis.Cmdable) limiter.Limiter {
	return limiter.NewRedisSlidingWindowLimiter(cmdable, time.Second, 1000)
}

func NewRedisIpRouteLimiter(cmdable redis.Cmdable) limiter.Limiter {
	return limiter.NewRedisSlidingWindowLimiter(cmdable, time.Second, 1000)
}
