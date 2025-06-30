//go:build wireinject

package main

import (
	"GoProj/wedy/interactive/events"
	"GoProj/wedy/interactive/grpc"
	"GoProj/wedy/interactive/ioc"
	"GoProj/wedy/interactive/repository"
	"GoProj/wedy/interactive/repository/cache"
	"GoProj/wedy/interactive/repository/dao"
	"GoProj/wedy/interactive/service"
	"github.com/google/wire"
)

var integrationSvcSet = wire.NewSet(
	cache.NewInteractiveRedisCache,
	repository.NewCachedInteractiveRepository,
	dao.NewGORMInteractiveDAO,
	service.NewInteractiveService,
)

var thirdPartySet = wire.NewSet(
	ioc.InitDB,
	ioc.InitSaramaClient,
	ioc.InitRedis,
	ioc.InitLogger)

func InitAPP() *App {
	wire.Build(thirdPartySet, integrationSvcSet, ioc.InitConsumer, events.NewInteractiveWatchEventConsumer, grpc.NewInteractiveServiceServer, ioc.NewGrpcxServer, wire.Struct(new(App), "*"))
	return new(App)
}
