package job

import (
	"GoProj/wedy/internal/service"
	"GoProj/wedy/pkg/logger"
	"context"
	rclock "github.com/gotomicro/redis-lock"
	"sync"
	"time"
)

type RankingJob struct {
	svc       service.RankingService
	timeout   time.Duration
	client    *rclock.Client
	key       string
	l         logger.LoggerV1
	lock      *rclock.Lock
	locallock *sync.Mutex
}

func NewRankingJob(svc service.RankingService, timeout time.Duration, client *rclock.Client, l logger.LoggerV1) *RankingJob {
	return &RankingJob{svc: svc, timeout: timeout, client: client, key: "job:ranking", l: l, locallock: &sync.Mutex{}}
}

func (r *RankingJob) Name() string {
	return "RankingJob"
}
func (r *RankingJob) Run() error {
	r.locallock.Lock()
	lock := r.lock
	if lock == nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
		defer cancel()
		locker, err := r.client.Lock(ctx, r.key, r.timeout, &rclock.FixIntervalRetry{
			Interval: time.Millisecond * 100,
			Max:      3,
		}, time.Second)
		if err != nil {
			r.locallock.Unlock()
			r.l.Warn("Lock failed", logger.Error(err))
			return nil
		}
		r.lock = locker
		r.locallock.Unlock()
		go func() {
			er := locker.AutoRefresh(r.timeout/2, r.timeout)
			if er != nil {
				r.locallock.Lock()
				r.lock = nil
				r.locallock.Unlock()
			}
		}()
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	return r.svc.TopN(ctx)
}

//func (r *RankingJob) Run() error {
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
//	defer cancel()
//	locker, err := r.client.Lock(ctx, r.key, r.timeout, &rclock.FixIntervalRetry{
//		Interval: time.Millisecond * 100,
//		Max:      3,
//	}, time.Second)
//	if err != nil {
//		return err
//	}
//	defer func() {
//		ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
//		defer cancel()
//		err := locker.Unlock(ctx)
//		if err != nil {
//			r.l.Error("ranking job Unlock failed", logger.Error(err))
//		}
//	}()
//	ctx, cancel = context.WithTimeout(context.Background(), r.timeout)
//	defer cancel()
//	return r.svc.TopN(ctx)
//}
