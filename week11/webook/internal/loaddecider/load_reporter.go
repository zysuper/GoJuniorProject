package loaddecider

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

// Reporter 将系统负载回报到某个存储里
type Reporter interface {
	// Report 报告某个服务的负载
	Report(context context.Context, id string, load int32)
}

type RedisZSetReporter struct {
	redis      redis.Cmdable
	expiration time.Duration
}

// NopReporter 用于测试.
type NopReporter struct {
}

func NewRedisZSetReporter(redis redis.Cmdable) Reporter {
	return &RedisZSetReporter{redis: redis, expiration: time.Minute}
}

func NewNopReporter() *NopReporter {
	return &NopReporter{}
}

func (r *RedisZSetReporter) Report(context context.Context, id string, load int32) {
	r.redis.ZAdd(context, LoadSetKey, redis.Z{Score: float64(load), Member: id})
	// 标记过期时间.
	r.redis.Set(context, LoadExpirePrefix+id, "", r.expiration)
}

func (n *NopReporter) Report(context context.Context, id string, load int32) {
}
