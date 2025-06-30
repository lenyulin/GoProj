package service

import (
	"GoProj/wedy/internal/domian"
	"GoProj/wedy/internal/repository"
	"GoProj/wedy/pkg/logger"
	"context"
	"time"
)

type CronJobService interface {
	Preempt(ctx context.Context) (domian.Job, error)
	ResetNextTime(ctx context.Context, j domian.Job) error
}
type cronJobService struct {
	repo            repository.JobRepository
	l               logger.LoggerV1
	refreshInterval time.Duration
}

func newCronJobService(repo repository.JobRepository, l logger.LoggerV1) CronJobService {
	return &cronJobService{repo: repo, l: l, refreshInterval: time.Minute}
}

func (c *cronJobService) ResetNextTime(ctx context.Context, j domian.Job) error {
	nextTime := j.NextTime()
	return c.repo.UpdateNextTime(ctx, j.Id, nextTime)
}

func (c *cronJobService) Preempt(ctx context.Context) (domian.Job, error) {
	j, err := c.repo.Preempt(ctx)
	if err != nil {
		return domian.Job{}, err
	}
	ticker := time.NewTicker(c.refreshInterval)
	go func() {
		for range ticker.C {
			c.refresh(j.Id)
		}
	}()
	j.Cancel = func() {
		ticker.Stop()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		err := c.repo.Release(ctx, j.Id)
		if err != nil {
			c.l.Error("Release job failed", logger.Error(err), logger.Int64("job id", j.Id))
		}
	}
	return j, nil
}
func (c *cronJobService) refresh(id int64) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err := c.repo.UpdateUtime(ctx, id)
	if err != nil {
		c.l.Error("refresh failed", logger.Error(err), logger.Int64("job id", id))
	}
}
