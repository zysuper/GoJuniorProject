package repository

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
)

type TopNRepository interface {
	GetTopLikeN(context context.Context, bizKey string, size int) ([]domain.TopLike, error)
}

type CachedTopNServiceRepository struct {
	dao   dao.LikeTopNDAO
	cache cache.LikeTopCache
}

func NewCachedTopNServiceRepository(dao dao.LikeTopNDAO, cache cache.LikeTopCache) TopNRepository {
	return &CachedTopNServiceRepository{dao: dao, cache: cache}
}

func (c *CachedTopNServiceRepository) GetTopLikeN(context context.Context, bizKey string, size int) ([]domain.TopLike, error) {
	l, error := c.cache.GetTopLikeN(context, bizKey, size)
	if error == nil {
		return l, nil
	}

	list, error := c.dao.QueryLikeNList(context, bizKey, size)

	if error != nil {
		return []domain.TopLike{}, error
	}

	// 异步保存缓存.
	go func() {
		c.cache.SaveTopLikeN(context.Background(), bizKey, list)
	}()

	return list, nil
}
