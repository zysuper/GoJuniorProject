package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func Test_msgDAO_Create(t *testing.T) {
	tests := []struct {
		name    string
		mock    func(t *testing.T) *sql.DB
		msg     Msg
		wantErr error
	}{
		{
			name: "创建成功",
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mockRes := sqlmock.NewResult(123, 1)
				// 这边要求传入的是 sql 的正则表达式
				mock.ExpectExec("INSERT INTO .*").
					WillReturnResult(mockRes)
				return db
			},
			msg: Msg{Id: 123, Args: "{}"},
		},
		{
			name: "创建失败",
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				// 这边要求传入的是 sql 的正则表达式
				mock.ExpectExec("INSERT INTO .*").
					WillReturnError(errors.New("数据库错误"))
				return db
			},
			msg:     Msg{Id: 123, Args: "{}"},
			wantErr: errors.New("数据库错误"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlDB := tt.mock(t)
			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      sqlDB,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				DisableAutomaticPing:   true,
				SkipDefaultTransaction: true,
			})
			assert.NoError(t, err)
			dao := NewMsgDao(db)
			_, err = dao.Create(context.Background(), tt.msg)

			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func Test_msgDAO_FindById(t *testing.T) {
	tests := []struct {
		name    string
		mock    func(t *testing.T) *sql.DB
		uid     int64
		wantErr error
		wantMsg Msg
	}{
		{
			name: "查询的到",
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				// 这边要求传入的是 sql 的正则表达式
				mock.ExpectQuery("SELECT (.+) FROM `msgs` .*").
					WithArgs(123).
					WillReturnRows(mock.NewRows([]string{"id", "args"}).AddRow(123, "{}"))
				return db
			},
			uid:     123,
			wantErr: nil,
			wantMsg: Msg{Id: 123, Args: "{}"},
		},
		{
			name: "查询不到",
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				// 这边要求传入的是 sql 的正则表达式
				mock.ExpectQuery("SELECT (.+) FROM `msgs` .*").
					WithArgs(123).WillReturnRows(mock.NewRows([]string{"id", "args"}))
				return db
			},
			uid:     123,
			wantErr: errors.New("record not found"),
			wantMsg: Msg{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlDB := tt.mock(t)
			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      sqlDB,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				DisableAutomaticPing:   true,
				SkipDefaultTransaction: true,
			})
			assert.NoError(t, err)
			dao := NewMsgDao(db)
			msg, err := dao.FindById(context.Background(), tt.uid)

			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantMsg, msg)
		})
	}
}
