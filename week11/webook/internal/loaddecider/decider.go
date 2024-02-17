package loaddecider

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"github.com/redis/go-redis/v9"
	"time"
)

// Decider 负载决策器.
type Decider interface {
	// IsVictory 是否由当前节点执行
	IsVictory(context context.Context) bool
}

// RedisZSetDecider redis zset 决策器.
type RedisZSetDecider struct {
	redis   redis.Cmdable
	l       logger.LoggerV1
	timeout time.Duration
	id      Identity
}

func NewRedisZSetDecider(redis redis.Cmdable, l logger.LoggerV1, id Identity) Decider {
	return &RedisZSetDecider{redis: redis, l: l, id: id, timeout: time.Minute * 1}
}

func (r *RedisZSetDecider) IsVictory(ctx context.Context) bool {
	result, err := r.redis.ZRangeWithScores(ctx, LoadSetKey, 0, 0).Result()

	if err != nil {
		// redis 没有记录，就自己干吧.
		r.l.Error("redis 查不到或者异常", logger.Error(err))
		return true
	}

	s, ok := result[0].Member.(string)

	if !ok {
		r.l.Error("member 不为预期类型， 谁干的?")
		// 异常了，只能降级自己干了.
		return true
	}

	exists, err := r.redis.Exists(ctx, LoadExpirePrefix+s).Result()
	if err != nil {
		r.l.Error("判断是否已经过期失败", logger.Error(err))
		// 异常了，只能降级自己干了.
		return true
	}

	if exists == 0 {
		// 这个值过期了，需要从 zset 移除.
		// 说明这个节点的服务很久没有更新自己的 load 了，挂了?
		cxt, cancel := context.WithTimeout(context.Background(), time.Second)
		// 有并发问题，最好用 redis lua 脚本干。
		_, err := r.redis.ZRem(cxt, LoadSetKey, s).Result()
		if err != nil {
			r.l.Error("zrem 失败", logger.Error(err))
		}
		cancel()
		// 这是过期数据，只能降级自己干了.
		return true
	}

	return r.id.Id() == s
}
