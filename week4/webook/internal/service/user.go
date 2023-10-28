package service

import (
	"context"
	"errors"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	"github.com/gin-gonic/gin"
)

var (
	ErrDuplicateUser         = repository.ErrDuplicateUser
	ErrInvalidUserOrPassword = errors.New("用户不存在或者密码不对")
)

type UserService interface {
	Signup(ctx context.Context, u domain.User) error
	Login(ctx context.Context, email string, password string) (domain.User, error)
	UpdateUserInfo(ctx context.Context, user domain.User) error
	FindById(ctx *gin.Context, uid int64) (domain.User, error)
	FindOrCreate(ctx *gin.Context, phone string) (domain.User, error)
}

type userService struct {
	repo repository.UserRepository
	pv   PasswordValidateService
}

func (svc *userService) FindOrCreate(ctx *gin.Context, phone string) (domain.User, error) {
	// 先找一下，我们认为，大部分用户是已经存在的用户
	u, err := svc.repo.FindByPhone(ctx, phone)
	if err != repository.ErrUserNotFound {
		// 有两种情况
		// err == nil, u 是可用的
		// err != nil，系统错误，
		return u, err
	}
	// 用户没找到
	err = svc.repo.Create(ctx, domain.User{
		Phone: phone,
	})
	// 有两种可能，一种是 err 恰好是唯一索引冲突（phone）
	// 一种是 err != nil，系统错误
	if err != nil && err != ErrDuplicateUser {
		return domain.User{}, err
	}
	// 要么 err ==nil，要么ErrDuplicateUser，也代表用户存在
	// 主从延迟，理论上来讲，强制走主库
	return svc.repo.FindByPhone(ctx, phone)
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
