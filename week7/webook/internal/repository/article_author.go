package repository

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
)

type ArticleAuthorRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
}
