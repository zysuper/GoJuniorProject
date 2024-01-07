package service

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository"
)

type TopNService interface {
	GetTopLikeN(context context.Context, bizKey string, size int) ([]domain.TopLike, error)
}

type DefaultTopNService struct {
	repository repository.TopNRepository
}

func NewDefaultTopNService(repository repository.TopNRepository) TopNService {
	return &DefaultTopNService{repository: repository}
}

func (d *DefaultTopNService) GetTopLikeN(context context.Context, bizKey string, size int) ([]domain.TopLike, error) {
	return d.repository.GetTopLikeN(context, bizKey, size)
}
