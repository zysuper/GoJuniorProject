// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package startup

import (
	"gitee.com/geekbang/basic-go/webook/follow/grpc"
	"gitee.com/geekbang/basic-go/webook/follow/repository"
	"gitee.com/geekbang/basic-go/webook/follow/repository/cache"
	"gitee.com/geekbang/basic-go/webook/follow/repository/dao"
	"gitee.com/geekbang/basic-go/webook/follow/service"
)

// Injectors from wire.go:

func InitServer() *grpc.FollowServiceServer {
	gormDB := InitTestDB()
	followRelationDao := dao.NewGORMFollowRelationDAO(gormDB)
	cmdable := InitRedis()
	followCache := cache.NewRedisFollowCache(cmdable)
	loggerV1 := InitLog()
	followRepository := repository.NewFollowRelationRepository(followRelationDao, followCache, loggerV1)
	followRelationService := service.NewFollowRelationService(followRepository)
	followServiceServer := grpc.NewFollowRelationServiceServer(followRelationService)
	return followServiceServer
}
