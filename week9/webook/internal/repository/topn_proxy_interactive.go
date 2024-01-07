package repository

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
)

// ProxyLikeTopCachedInteractiveRepository 增加了对 like top 的维护.
type ProxyLikeTopCachedInteractiveRepository struct {
	interactiveRepo InteractiveRepository
	likeTopCache    cache.LikeTopCache
}

func NewProxyLikeTopCachedInteractiveRepository(dao dao.InteractiveDAO, cache cache.InteractiveCache, likeTopCache cache.LikeTopCache) InteractiveRepository {
	return &ProxyLikeTopCachedInteractiveRepository{interactiveRepo: NewCachedInteractiveRepository(dao, cache), likeTopCache: likeTopCache}
}

func (p ProxyLikeTopCachedInteractiveRepository) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	return p.interactiveRepo.IncrReadCnt(ctx, biz, bizId)
}

func (p ProxyLikeTopCachedInteractiveRepository) IncrLike(ctx context.Context, biz string, id int64, uid int64) error {
	go func() {
		p.likeTopCache.IncLikeCnt(context.Background(), biz, id)
	}()
	return p.interactiveRepo.IncrLike(ctx, biz, id, uid)
}

func (p ProxyLikeTopCachedInteractiveRepository) DecrLike(ctx context.Context, biz string, id int64, uid int64) error {
	go func() {
		p.likeTopCache.DecLikeCnt(context.Background(), biz, id)
	}()
	return p.interactiveRepo.DecrLike(ctx, biz, id, uid)
}

func (p ProxyLikeTopCachedInteractiveRepository) AddCollectionItem(ctx context.Context, biz string, id int64, cid int64, uid int64) error {
	return p.interactiveRepo.AddCollectionItem(ctx, biz, id, cid, uid)
}
