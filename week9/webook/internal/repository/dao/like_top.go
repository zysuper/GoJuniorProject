package dao

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gorm.io/gorm"
)

type LikeTopNDAO interface {
	QueryLikeNList(context context.Context, bizKey string, n int) []domain.TopLike
}

type GORMLikeTopNDAO struct {
	Db *gorm.DB
}

func NewGORMLikeTopNDAO(db *gorm.DB) *GORMLikeTopNDAO {
	return &GORMLikeTopNDAO{Db: db}
}

func (g *GORMLikeTopNDAO) QueryLikeNList(context context.Context, bizKey string, n int) []domain.TopLike {
	var results []domain.TopLike
	g.Db.
		WithContext(context).
		Select("a.id aid,a.title title,a.author_id as author_id,u.nickname,i.like_cnt as like_count").
		Table("interactives i").
		Joins("join articles a on a.id = i.biz_id join users u on a.author_id = u.id").
		Where(`i.biz = ?`, bizKey).
		Order("i.like_cnt desc").
		Limit(n).
		Scan(&results)
	return results
}
