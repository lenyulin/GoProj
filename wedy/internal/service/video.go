package service

import (
	"GoProj/wedy/internal/domian"
	"GoProj/wedy/internal/events/video"
	"GoProj/wedy/internal/repository"
	"GoProj/wedy/pkg/logger"
	"context"
	"time"
)

type VideoService interface {
	Publish(ctx context.Context, video domian.Video) error
	Update(ctx context.Context, video domian.Video) error
	Withdrawn(ctx context.Context, video domian.Video) error
	GetByAuthor(ctx context.Context, id int64, limit int, offset int) ([]domian.Video, error)
	GetById(ctx context.Context, id, uid int64) (domian.Video, error)
	ListPub(ctx context.Context, start time.Time, offset int, limit int) ([]domian.Video, error)
}

type videoService struct {
	repo     repository.VideoRepository
	producer video.VideoProducer
	log      logger.LoggerV1
}

func NewVideoService(repo repository.VideoRepository, producer video.VideoProducer, log logger.LoggerV1) VideoService {
	return &videoService{
		repo:     repo,
		producer: producer,
		log:      log,
	}
}
func (svc *videoService) ListPub(ctx context.Context, start time.Time, offset int, limit int) ([]domian.Video, error) {
	return svc.repo.ListPub(ctx, start, offset, limit)
}
func (svc *videoService) GetByAuthor(ctx context.Context, id int64, limit int, offset int) ([]domian.Video, error) {
	return svc.repo.GetByAuthor(ctx, id, limit, offset)
}
func (svc *videoService) GetById(ctx context.Context, id, uid int64) (domian.Video, error) {
	res, err := svc.repo.GetById(ctx, id)
	if err == nil {
		go func() {
			er := svc.producer.VideoProduceWatchEvent(video.WatchEvent{
				Vid: id,
				Uid: uid,
			})
			if er != nil {
				svc.log.Error("Send to watch event failed", logger.Error(er))
			}
		}()
	}
	return res, err
}
func (svc *videoService) Update(ctx context.Context, video domian.Video) error {
	return svc.repo.Update(ctx, video)
}
func (svc *videoService) Publish(ctx context.Context, video domian.Video) error {
	video.Status = domian.VideoStatusPublished
	return svc.repo.Create(ctx, video)
}
func (svc *videoService) Withdrawn(ctx context.Context, video domian.Video) error {
	video.Status = domian.VideoStatusUnpublished
	return svc.repo.SyncStatus(ctx, video)
}
