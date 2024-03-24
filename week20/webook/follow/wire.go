//go:build wireinject

package main

import (
	"gitee.com/geekbang/basic-go/webook/follow/events"
	grpc2 "gitee.com/geekbang/basic-go/webook/follow/grpc"
	"gitee.com/geekbang/basic-go/webook/follow/ioc"
	"gitee.com/geekbang/basic-go/webook/follow/repository"
	"gitee.com/geekbang/basic-go/webook/follow/repository/cache"
	"gitee.com/geekbang/basic-go/webook/follow/repository/dao"
	"gitee.com/geekbang/basic-go/webook/follow/service"
	"github.com/google/wire"
)

var serviceProviderSet = wire.NewSet(
	dao.NewGORMFollowRelationDAO,
	cache.NewRedisFollowCache,
	repository.NewFollowRelationRepository,
	service.NewFollowRelationService,
	grpc2.NewFollowRelationServiceServer,
)

var thirdProvider = wire.NewSet(
	ioc.InitDB,
	ioc.InitLogger,
	ioc.InitRedis,
	ioc.InitSaramaClient,
	ioc.InitEtcdClient,
)

func Init() *App {
	wire.Build(
		thirdProvider,
		serviceProviderSet,
		ioc.InitGRPCxServer,
		events.NewFollowConsumer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
