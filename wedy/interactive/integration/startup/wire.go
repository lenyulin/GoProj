//go:build wireinject

package startup

import (
	"GoProj/wedy/interactive/grpc"
	"GoProj/wedy/interactive/repository"
	"GoProj/wedy/interactive/repository/cache"
	"GoProj/wedy/interactive/repository/dao"
	"GoProj/wedy/interactive/service"
	"github.com/google/wire"
)

var thirdParty = wire.NewSet(InitDB, InitRedis, InitLogger)
var integrationSvcSet = wire.NewSet(
	cache.NewInteractiveRedisCache,
	repository.NewCachedInteractiveRepository,
	dao.NewGORMInteractiveDAO,
	service.NewInteractiveService,
)

func InitIntegrationService() *grpc.InteractiveServiceServer {
	wire.Build(thirdParty, integrationSvcSet, grpc.NewInteractiveServiceServer)
	return new(grpc.InteractiveServiceServer)
}
