//go:build wireinject

package main

import (
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	"gitee.com/geekbang/basic-go/webook/internal/web"
	"gitee.com/geekbang/basic-go/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		ioc.InitDB, ioc.InitRedis,
		cache.NewRedisUserCache, cache.NewRedisCodeCache,
		dao.NewUserDAO,
		repository.NewUserRepository, repository.NewCodeRepository,
		ioc.InitSms,
		service.NewCodeService,
		service.NewPasswordValidator,
		service.NewUserService,
		web.NewUserHandler,
		ioc.InitGinMiddlewares,
		ioc.InitWebServer)
	return gin.Default()
}
