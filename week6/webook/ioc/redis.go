package ioc

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedis() redis.Cmdable {
	// 第三方资源.
	return redis.NewClient(&redis.Options{
		Addr: viper.GetString("redis.addr"),
	})
}
