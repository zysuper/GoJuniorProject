package repository

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache/code"
)

var ErrCodeVerifyTooMany = code.ErrCodeVerifyTooMany
var ErrCodeSendTooMany = code.ErrCodeSendTooMany

type CodeRepository interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}

type CachedCodeRepository struct {
	cache code.CodeCache
}

func (c *CachedCodeRepository) Set(ctx context.Context, biz, phone, code string) error {
	return c.cache.Set(ctx, biz, phone, code)
}

func (c *CachedCodeRepository) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	return c.cache.Verify(ctx, biz, phone, code)
}

func NewCodeRepository(cache code.CodeCache) CodeRepository {
	return &CachedCodeRepository{
		cache: cache,
	}
}
