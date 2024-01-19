package dao

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestArticleGORMDAO_ListPub(t *testing.T) {
	db, _ := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"), &gorm.Config{})
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		ctx    context.Context
		start  time.Time
		offset int
		limit  int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []PublishedArticle
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:   "第一个测试",
			fields: fields{db: db},
			args:   args{ctx: context.Background(), start: time.Now().Add(-14 * time.Hour * 24), offset: 0, limit: 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &ArticleGORMDAO{
				db: tt.fields.db,
			}
			got, err := a.ListPub(tt.args.ctx, tt.args.start, tt.args.offset, tt.args.limit)
			if !tt.wantErr(t, err, fmt.Sprintf("ListPub(%v, %v, %v, %v)", tt.args.ctx, tt.args.start, tt.args.offset, tt.args.limit)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ListPub(%v, %v, %v, %v)", tt.args.ctx, tt.args.start, tt.args.offset, tt.args.limit)
		})
	}
}
