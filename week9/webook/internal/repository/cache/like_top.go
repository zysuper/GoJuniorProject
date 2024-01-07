package cache

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"github.com/redis/go-redis/v9"
	"strconv"
)

type LikeTopCache interface {
	// GetTopLikeN 获取排名前 n 的数据
	GetTopLikeN(context context.Context, bizKey string, n int) ([]domain.TopLike, error)
	// SaveTopLikeN 首次存储数据库查询的排名前 n 的数据
	SaveTopLikeN(context context.Context, bizKey string, l []domain.TopLike) error
	// IncLikeCnt 排名前 n 的数据又有新的👍
	IncLikeCnt(context context.Context, bizKey string, aid int64) error
	// DecLikeCnt 排名前 n 的数据又有新的取消👍
	DecLikeCnt(context context.Context, bizKey string, aid int64) error
	// MemberExists 看看当前被👍的小伙子是否在当前排名里面
	MemberExists(context context.Context, bizKey string, aid int64) (bool, error)
}

type RedisLikeTopCache struct {
	client redis.Cmdable
}

func (r *RedisLikeTopCache) MemberExists(context context.Context, bizKey string, aid int64) (bool, error) {
	_, err := r.client.ZScore(context, keyWithSuffix(bizKey), strconv.FormatInt(aid, 10)).Result()
	if err == redis.Nil {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *RedisLikeTopCache) IncLikeCnt(context context.Context, bizKey string, aid int64) error {
	// 存在才操作，否则，忽略.
	if s, _ := r.MemberExists(context, bizKey, aid); s == true {
		return r.client.ZIncrBy(context, keyWithSuffix(bizKey), 1, strconv.FormatInt(aid, 10)).Err()
	}
	return nil
}

func (r *RedisLikeTopCache) DecLikeCnt(context context.Context, bizKey string, aid int64) error {
	// 存在才操作，否则，忽略.
	if s, _ := r.MemberExists(context, bizKey, aid); s == true {
		return r.client.ZIncrBy(context, keyWithSuffix(bizKey), -1, strconv.FormatInt(aid, 10)).Err()
	}
	return nil
}

func (r *RedisLikeTopCache) AddNewLikeN(context context.Context, bizKey string, like domain.TopLike) error {
	return r.client.ZAdd(context, keyWithSuffix(bizKey), redis.Z{Score: float64(like.LikeCount), Member: like.Aid}).Err()
}

func (r *RedisLikeTopCache) GetTopLikeN(context context.Context, bizKey string, n int) ([]domain.TopLike, error) {
	result, err := r.client.ZRevRangeWithScores(context, keyWithSuffix(bizKey), 0, int64(n)).Result()
	if err != nil {
		return nil, err
	}

	var ret []domain.TopLike
	for _, z := range result {
		i, _ := strconv.ParseInt(z.Member.(string), 10, 64)
		ret = append(ret, domain.TopLike{LikeCount: int64(z.Score), Aid: i})
	}

	return ret, nil
}

func (r *RedisLikeTopCache) SaveTopLikeN(context context.Context, bizKey string, l []domain.TopLike) error {
	var members []redis.Z
	for _, like := range l {
		members = append(members, redis.Z{Score: float64(like.LikeCount), Member: like.Aid})
	}
	return r.client.ZAdd(context, keyWithSuffix(bizKey), members...).Err()
}

func NewRedisLikeTopCache(client redis.Cmdable) LikeTopCache {
	return &RedisLikeTopCache{client: client}
}

func keyWithSuffix(key string) string {
	return key + "_topn"
}
