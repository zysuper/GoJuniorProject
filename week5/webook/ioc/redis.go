package ioc

import (
	"gitee.com/geekbang/basic-go/webook/config"
	"github.com/redis/go-redis/v9"
)

func InitRedis() redis.Cmdable {
	// 第三方资源.
	return redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
}
