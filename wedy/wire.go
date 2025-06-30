//go:build wireinject

package main

import (
	repository3 "GoProj/wedy/comment/repository"
	dao4 "GoProj/wedy/comment/repository/dao"
	service3 "GoProj/wedy/comment/service"
	"GoProj/wedy/interactive/events"
	repository2 "GoProj/wedy/interactive/repository"
	cache2 "GoProj/wedy/interactive/repository/cache"
	dao3 "GoProj/wedy/interactive/repository/dao"
	service2 "GoProj/wedy/interactive/service"
	"GoProj/wedy/internal/events/video"
	"GoProj/wedy/internal/repository"
	"GoProj/wedy/internal/repository/cache"
	"GoProj/wedy/internal/repository/dao"
	dao2 "GoProj/wedy/internal/repository/dao/video"
	"GoProj/wedy/internal/service"
	"GoProj/wedy/internal/web"
	"GoProj/wedy/ioc"
	"GoProj/wedy/pkg/logger"
	"GoProj/wedy/pkg/oss"
	"github.com/IBM/sarama"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func InitializeVideoRepository(db *gorm.DB, ossUploader oss.OSSHandler, cache cache.VideoCache) repository.VideoRepository {
	gormDAO := dao2.NewGORMVideoDAO(db, cache)
	vdao := dao2.NewOSSVideoDao(ossUploader, gormDAO)
	return repository.NewVideoRepository(vdao, cache)
}
func InitOssHdl() oss.OSSHandler {
	return oss.NewOSSHandler(ioc.InitOSS())
}

func InitCachedInteractiveRepository(dao dao3.InteractiveDao, cache cache2.InteractiveCache) repository2.InteractiveRepository {
	return repository2.NewCachedInteractiveRepository(dao, cache)
}
func InitCommentRepository(db *gorm.DB) repository3.CommentRepository {
	dao := dao4.NewGORMCommentDAO(db)
	return repository3.NewCachedCommentRepository(dao)
}
func InitInteractiveDAO(db *gorm.DB) dao3.InteractiveDao {
	return dao3.NewGORMInteractiveDAO(db)
}
func InitInteractiveCache(client redis.Cmdable) cache2.InteractiveCache {
	return cache2.NewInteractiveRedisCache(client)
}
func InitInteractiveWatchEventConsumer(repo repository2.InteractiveRepository, client sarama.Client, log logger.LoggerV1) *events.InteractiveWatchEventConsumer {
	return events.NewInteractiveWatchEventConsumer(repo, client, log)
}

var rankingSvcSet = wire.NewSet(cache.NewRedisRankingCache, repository.NewCacheRankingRepository, service.NewBatchRankingService)

func InitWebServer() *App {
	wire.Build(
		ioc.InitDB, InitializeVideoRepository, InitOssHdl, ioc.InitLogger, InitCommentRepository, ioc.InitSaramaClient, ioc.InitSyncProducer,
		ioc.InitRedis, ioc.InitConsumer, ioc.InitWebServer, ioc.InitMiddleware,
		rankingSvcSet, ioc.InitJobs, ioc.InitRankingJob, ioc.InitRLockClent,
		dao.NewUserDAO, InitInteractiveDAO, InitInteractiveCache, service2.NewInteractiveService, InitCachedInteractiveRepository,
		cache.NewCodeCache, cache.NewUserCache, cache.NewVideoRedisCache,

		repository.NewCachedUserRepository, repository.NewCodeRepository, InitInteractiveWatchEventConsumer,
		ioc.InitSMSService, video.NewSaramaSyncProducer,

		service.NewUserService, service.NewCodeService, service.NewVideoService, service3.NewCommentService,
		web.NewUserHandler, web.NewVideoHandler, web.NewCommentHandler,

		wire.Struct(new(App), "*"),
	)
	return new(App)
}

//func InitWebServer() *App {
//	wire.Build(
//		ioc.InitDB, InitializeVideoRepository, InitOssHdl, interactiveSvc,
//		ioc.InitRedis,
//		ioc.InitMongoDB,
//		dao.NewUserDAO, dao.NewGORMCommentDAO, dao2.NewGORMVideoDAO, dao2.NewOSSVideoDao,
//		cache.NewCodeCache, cache.NewUserCache, cache.NewVideoRedisCache,
//		repository.NewCachedUserRepository, repository.NewCodeRepository,
//		repository.NewCachedCommentRepository,
//		ioc.InitSMSService,
//		service.NewUserService, service.NewCodeService, service.NewVideoService, service.NewCommentService,
//
//		web.NewUserHandler, web.NewVideoHandler, web.NewCommentHandler,
//		ioc.InitWebServer, ioc.InitMiddleware, dao2.NewMongoDBDAO,
//
//		wire.Struct(new(App), "*"),
//	)
//	return App
//}
