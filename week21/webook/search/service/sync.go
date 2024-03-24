package service

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/search/domain"
	"gitee.com/geekbang/basic-go/webook/search/repository"
)

type SyncService interface {
	InputArticle(ctx context.Context, article domain.Article) error
	InputUser(ctx context.Context, user domain.User) error
	InputAny(ctx context.Context, idxName, docID, data string) error
}

type syncService struct {
	userRepo    repository.UserRepository
	articleRepo repository.ArticleRepository
	anyRepo     repository.AnyRepository
}

func (s *syncService) InputAny(ctx context.Context, index, docID, data string) error {
	return s.anyRepo.Input(ctx, index, docID, data)
}

func (s *syncService) InputArticle(ctx context.Context, article domain.Article) error {
	return s.articleRepo.InputArticle(ctx, article)
}

func (s *syncService) InputUser(ctx context.Context, user domain.User) error {
	return s.userRepo.InputUser(ctx, user)
}

func NewSyncService(
	anyRepo repository.AnyRepository,
	userRepo repository.UserRepository,
	articleRepo repository.ArticleRepository) SyncService {
	return &syncService{
		userRepo:    userRepo,
		articleRepo: articleRepo,
		anyRepo:     anyRepo,
	}
}
