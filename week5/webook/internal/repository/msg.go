package repository

import (
	"context"
	"encoding/json"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
	"log"
)

type MsgRepository interface {
	Create(ctx context.Context, msg domain.Msg) (int64, error)
	FindById(ctx context.Context, uid int64) (domain.Msg, error)
}

type msgRepository struct {
	dao dao.MsgDAO
}

func NewMsgRepository(dao dao.MsgDAO) MsgRepository {
	return &msgRepository{dao: dao}
}

func (m *msgRepository) Create(ctx context.Context, msg domain.Msg) (int64, error) {
	entity, err := toEntity(msg)
	if err != nil {
		return 0, err
	}
	return m.dao.Create(ctx, entity)
}

func toEntity(msg domain.Msg) (dao.Msg, error) {
	json, err := json.Marshal(msg)
	if err != nil {
		log.Println("marshal 失败.")
		return dao.Msg{}, err
	}
	return dao.Msg{
		Args: string(json),
	}, nil
}

func (m *msgRepository) FindById(ctx context.Context, uid int64) (domain.Msg, error) {
	msg, err := m.dao.FindById(ctx, uid)
	if err != nil {
		return domain.Msg{}, err
	}

	return toDomain(msg), err
}

func toDomain(msg dao.Msg) domain.Msg {
	var rm domain.Msg
	err := json.Unmarshal([]byte(msg.Args), &rm)
	if err != nil {
		log.Println(err)
		return domain.Msg{}
	}
	return rm
}
