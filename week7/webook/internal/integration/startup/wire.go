//go:build wireinject

package startup

import (
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache/code"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	"gitee.com/geekbang/basic-go/webook/internal/web"
	ijwt "gitee.com/geekbang/basic-go/webook/internal/web/jwt"
	"gitee.com/geekbang/basic-go/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var thirdPartySet = wire.NewSet( // 第三方依赖
	InitRedis, InitDB,
	InitLogger)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 第三方依赖
		thirdPartySet,
		// DAO 部分
		dao.NewUserDAO,
		dao.NewArticleGORMDAO,
		dao.NewMsgDao,

		// cache 部分
		code.NewRedisCodeCache, cache.NewRedisUserCache,

		// repository 部分
		repository.NewUserRepository,
		repository.NewCodeRepository,
		repository.NewCachedArticleRepository,
		repository.NewMsgRepository,

		// Service 部分
		ioc.NewSmsCircuitBreaker,
		ioc.InitSms,
		service.NewUserService,
		service.NewCodeService,
		service.NewPasswordValidator,
		service.NewArticleService,
		ioc.InitWechatService,

		// handler 部分
		web.NewUserHandler,
		web.NewArticleHandler,
		web.NewOAuth2WechatHandler,
		ijwt.NewRedisJWTHandler,

		ioc.InitGinMiddlewares,
		ioc.InitWebServer,
	)
	return gin.Default()
}

func InitArticleHandler() *web.ArticleHandler {
	wire.Build(
		thirdPartySet,
		dao.NewArticleGORMDAO,
		service.NewArticleService,
		web.NewArticleHandler,
		repository.NewCachedArticleRepository)
	return &web.ArticleHandler{}
}
