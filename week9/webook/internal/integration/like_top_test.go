package integration

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
	"gitee.com/geekbang/basic-go/webook/ioc"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGORMLikeTopNDAO_queryLikeNList(t *testing.T) {
	g := &dao.GORMLikeTopNDAO{
		Db: ioc.InitDB(ioc.InitLogger()),
	}
	list, _ := g.QueryLikeNList(context.Background(), "article", 100)
	assert.NotEmpty(t, list)
}
