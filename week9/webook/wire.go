//go:build wireinject

package main

import (
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	"gitee.com/geekbang/basic-go/webook/internal/web"
	ijwt "gitee.com/geekbang/basic-go/webook/internal/web/jwt"
	"gitee.com/geekbang/basic-go/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var interactiveSvcSet = wire.NewSet(dao.NewGORMInteractiveDAO,
	cache.NewInteractiveRedisCache,
	repository.NewProxyLikeTopCachedInteractiveRepository,
	service.NewInteractiveService,
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 第三方依赖
		ioc.InitRedis, ioc.InitDB,
		ioc.InitLogger,
		// DAO 部分
		dao.NewUserDAO,
		dao.NewArticleGORMDAO,
		dao.NewGORMLikeTopNDAO, // 1. like top dao

		interactiveSvcSet,

		// cache 部分
		cache.NewCodeCache, cache.NewUserCache,
		cache.NewArticleRedisCache,
		cache.NewRedisLikeTopCache, // 2. like top cache

		// repository 部分
		repository.NewCachedUserRepository,
		repository.NewCodeRepository,
		repository.NewCachedArticleRepository,
		repository.NewCachedTopNServiceRepository, // 3. like top repository

		// Service 部分
		ioc.InitSMSService,
		ioc.InitWechatService,
		service.NewUserService,
		service.NewCodeService,
		service.NewArticleService,
		service.NewDefaultTopNService, // 4. like top n service

		// handler 部分
		web.NewUserHandler,
		web.NewArticleHandler,
		ijwt.NewRedisJWTHandler,
		web.NewOAuth2WechatHandler,
		web.NewLikeTopHandler, // 5. like top n handler
		ioc.InitGinMiddlewares,
		ioc.InitWebServer,
	)
	return gin.Default()
}
