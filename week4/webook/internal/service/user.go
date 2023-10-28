package service

import (
	"context"
	"errors"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	"github.com/gin-gonic/gin"
)

var (
	ErrDuplicateEmail        = repository.ErrDuplicateEmail
	ErrInvalidUserOrPassword = errors.New("用户不存在或者密码不对")
)

type UserService interface {
	Signup(ctx context.Context, u domain.User) error
	Login(ctx context.Context, email string, password string) (domain.User, error)
	UpdateUserInfo(ctx context.Context, user domain.User) error
	FindById(ctx *gin.Context, uid int64) (domain.User, error)
}

type userService struct {
	repo repository.UserRepository
	pv   PasswordValidateService
}

func NewUserService(repo repository.UserRepository, pv PasswordValidateService) UserService {
	return &userService{
		repo: repo,
		pv:   pv,
	}
}

func (svc *userService) Signup(ctx context.Context, u domain.User) error {
	hash, err := svc.pv.Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}

func (svc *userService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	// 检查密码对不对
	// err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	err = svc.pv.ComparePassword(u.Password, password)
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}

func (svc *userService) UpdateUserInfo(ctx context.Context,
	user domain.User) error {
	return svc.repo.UpdateUserById(ctx, user)
}

func (svc *userService) FindById(ctx *gin.Context, uid int64) (domain.User, error) {
	return svc.repo.FindById(ctx, uid)
}
