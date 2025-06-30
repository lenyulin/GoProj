package service

import (
	service2 "GoProj/wedy/interactive/service"
	"GoProj/wedy/internal/domian"
	"GoProj/wedy/internal/repository"
	"context"
	"errors"
	"github.com/ecodeclub/ekit/queue"
	"github.com/ecodeclub/ekit/slice"
	"math"
	"time"
)

type RankingService interface {
	TopN(ctx context.Context) error
	GetTopN(ctx context.Context) ([]domian.Video, error)
}
type BatchRankingService struct {
	interactiveSvc service2.InteractiveService
	videoSvc       VideoService
	batchSize      int
	scoreFun       func(readCnt int64, utime time.Time) float64
	n              int
	repo           repository.RankingRepository
}

func NewBatchRankingService(interactiveSvc service2.InteractiveService, videoSvc VideoService) RankingService {
	return &BatchRankingService{
		interactiveSvc: interactiveSvc,
		videoSvc:       videoSvc,
		batchSize:      2,
		scoreFun: func(readCnt int64, utime time.Time) float64 {
			duration := time.Since(utime).Seconds()
			return float64(readCnt-1) / math.Pow(duration+2, 1.5)
		},
		n: 3,
	}
}
func (b *BatchRankingService) GetTopN(ctx context.Context) ([]domian.Video, error) {
	return b.repo.GetRedisTopN(ctx)
}
func (b *BatchRankingService) TopN(ctx context.Context) error {
	videos, err := b.topN(ctx)
	if err != nil {
		return err
	}
	return b.repo.ReplaceRedisTopN(ctx, videos)
}
func (b *BatchRankingService) topN(ctx context.Context) ([]domian.Video, error) {
	type Score struct {
		score float64
		video domian.Video
	}
	topN := queue.NewPriorityQueue[Score](b.n,
		func(src Score, dst Score) int {
			if src.score > dst.score {
				return 1
			} else if src.score == dst.score {
				return 0
			} else {
				return -1
			}
		})
	offset := 0
	start := time.Now()
	ddl := start.Add(-30 * 24 * time.Hour)
	for {
		videos, err := b.videoSvc.ListPub(ctx, start, offset, b.batchSize)
		if err != nil {
			return []domian.Video{}, err
		}
		if len(videos) == 0 {
			break
		}
		ids := slice.Map(videos, func(idx int, src domian.Video) int64 {
			return src.Uid
		})
		intrMap, err := b.interactiveSvc.GetByIds(ctx, "video", ids)
		if err != nil {
			return []domian.Video{}, err
		}
		for _, v := range videos {
			intr := intrMap[v.Uid]
			score := b.scoreFun(intr.ReadCnt, v.Utime)
			ele := Score{score: score, video: v}
			err := topN.Enqueue(ele)
			if errors.Is(err, queue.ErrOutOfCapacity) {
				minEle, _ := topN.Dequeue()
				if minEle.score < score {
					_ = topN.Enqueue(ele)
				} else {
					_ = topN.Enqueue(minEle)
				}
			}
		}
		offset += len(videos)
		if len(videos) < b.batchSize || videos[len(videos)-1].Utime.Before(ddl) {
			break
		}
	}
	res := make([]domian.Video, topN.Len())
	for i := topN.Len() - 1; i >= 0; i-- {
		ele, _ := topN.Dequeue()
		res[i] = ele.video
	}
	return res, nil
}
