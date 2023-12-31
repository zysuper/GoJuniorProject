// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

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
)

// Injectors from wire.go:

func InitWebServer() *gin.Engine {
	cmdable := ioc.InitRedis()
	v := ioc.InitGinMiddlewares(cmdable)
	db := ioc.InitDB()
	userDAO := dao.NewUserDAO(db)
	userCache := cache.NewRedisUserCache(cmdable)
	userRepository := repository.NewUserRepository(userDAO, userCache)
	passwordValidateService := service.NewPasswordValidator()
	userService := service.NewUserService(userRepository, passwordValidateService)
	freecacheCache := ioc.InitLocalCache()
	codeCache := code.NewMemCodeCache(freecacheCache)
	codeRepository := repository.NewCodeRepository(codeCache)
	circuitBreaker := ioc.NewSmsCircuitBreaker()
	msgDAO := dao.NewMsgDao(db)
	msgRepository := repository.NewMsgRepository(msgDAO)
	smsService := ioc.InitSms(circuitBreaker, cmdable, msgRepository)
	codeService := service.NewCodeService(codeRepository, smsService)
	userHandler := web.NewUserHandler(userService, codeService)
	engine := ioc.InitWebServer(v, userHandler)
	return engine
}
