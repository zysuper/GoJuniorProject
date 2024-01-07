package dao

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gorm.io/gorm"
)

type LikeTopNDAO interface {
	QueryLikeNList(context context.Context, bizKey string, n int) ([]domain.TopLike, error)
}

type GORMLikeTopNDAO struct {
	Db *gorm.DB
}

func NewGORMLikeTopNDAO(db *gorm.DB) LikeTopNDAO {
	return &GORMLikeTopNDAO{Db: db}
}

func (g *GORMLikeTopNDAO) QueryLikeNList(context context.Context, bizKey string, n int) ([]domain.TopLike, error) {
	var results []domain.TopLike
	error := g.Db.
		WithContext(context).
		Select("i.biz_id as aid, i.like_cnt as like_count").
		Table("interactives i").
		Where(`i.biz = ?`, bizKey).
		Order("i.like_cnt desc").
		Limit(n).
		Scan(&results).
		Error
	return results, error
}
