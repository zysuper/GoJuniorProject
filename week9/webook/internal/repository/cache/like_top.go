package cache

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"github.com/redis/go-redis/v9"
)

type LikeTopCache interface {
	GetTopLikeN(context context.Context, bizKey string, n int) ([]domain.TopLike, error)
	SaveTopLikeN(context context.Context, bizKey string, l []domain.TopLike) error
	UpdateLike(context context.Context, bizKey string, like domain.TopLike) error
}

type RedisLikeTopCache struct {
	client redis.Cmdable
}

func (r *RedisLikeTopCache) UpdateLike(context context.Context, bizKey string, like domain.TopLike) error {
	return r.client.ZAdd(context, bizKey, redis.Z{Score: float64(like.LikeCount), Member: like.Aid}).Err()
}

func (r *RedisLikeTopCache) GetTopLikeN(context context.Context, bizKey string, n int) ([]domain.TopLike, error) {
	result, err := r.client.ZRevRangeWithScores(context, bizKey, 0, int64(n)).Result()
	if err != nil {
		return nil, err
	}

	var ret []domain.TopLike
	for _, z := range result {
		ret = append(ret, domain.TopLike{LikeCount: int64(z.Score), Aid: z.Member.(int64)})
	}

	return ret, nil
}

func (r *RedisLikeTopCache) SaveTopLikeN(context context.Context, bizKey string, l []domain.TopLike) error {
	var members []redis.Z
	for _, like := range l {
		members = append(members, redis.Z{Score: float64(like.LikeCount), Member: like.Aid})
	}
	return r.client.ZAdd(context, bizKey, members...).Err()
}

func NewRedisLikeTopCache(client redis.Cmdable) LikeTopCache {
	return &RedisLikeTopCache{client: client}
}
