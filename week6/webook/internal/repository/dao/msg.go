package dao

import (
	"context"
	"gorm.io/gorm"
)

type MsgDAO interface {
	Create(ctx context.Context, msg Msg) (int64, error)
	FindById(ctx context.Context, uid int64) (Msg, error)
}

type msgDAO struct {
	db *gorm.DB
}

func NewMsgDao(db *gorm.DB) MsgDAO {
	return &msgDAO{db: db}
}

func (m *msgDAO) Create(ctx context.Context, msg Msg) (int64, error) {
	err := m.db.WithContext(ctx).Create(&msg).Error
	if err != nil {
		return 0, err
	}
	return msg.Id, err
}

func (m *msgDAO) FindById(ctx context.Context, uid int64) (Msg, error) {
	var msg Msg
	err := m.db.WithContext(ctx).Where("id = ?", uid).First(&msg).Error
	return msg, err
}

type Msg struct {
	Id   int64  `gorm:"primaryKey,autoIncrement"`
	Args string `gorm:"not null"`
}
