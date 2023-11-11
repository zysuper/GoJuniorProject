package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

var ErrKeyNotExist = redis.Nil

type UserCache interface {
	Key(id int64) string
	Set(ctx context.Context, u domain.User) error
	Get(ctx context.Context, id int64) (domain.User, error)
}

type RedisUserCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func (r *RedisUserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := r.Key(id)
	result, err := r.cmd.Get(ctx, key).Result()

	if err != nil {
		return domain.User{}, err
	}

	var u domain.User

	json.Unmarshal([]byte(result), &u)

	return u, nil
}

func (r *RedisUserCache) Key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}

func (r *RedisUserCache) Set(ctx context.Context, u domain.User) error {
	uj, err := json.Marshal(u)
	if err != nil {
		return err
	}
	key := r.Key(u.Id)
	return r.cmd.Set(ctx, key, uj, r.expiration).Err()
}

func NewRedisUserCache(cmd redis.Cmdable) UserCache {
	return &RedisUserCache{
		cmd:        cmd,
		expiration: time.Minute * 15,
	}
}
