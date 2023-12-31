//go:build wireinject

package main

import (
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache/code"
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
		cache.NewRedisUserCache,
		// redis 实现的 code cache.
		// code.NewRedisCodeCache,
		// 内存本地实现的 code cache...
		ioc.InitLocalCache, code.NewMemCodeCache,
		dao.NewUserDAO,
		repository.NewUserRepository,
		repository.NewCodeRepository,
		ioc.InitSms,
		service.NewCodeService,
		service.NewPasswordValidator,
		service.NewUserService,
		web.NewUserHandler,
		ioc.InitGinMiddlewares,
		ioc.InitWebServer)
	return gin.Default()
}
