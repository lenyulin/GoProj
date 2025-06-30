package cache

import (
	"GoProj/wedy/internal/domian"
	"context"
	"errors"
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"time"
)

type RankingLocalCache struct {
	topN       *atomicx.Value[[]domian.Video]
	ddl        *atomicx.Value[time.Time]
	expiration time.Duration
}

func (r *RankingLocalCache) Set(ctx context.Context, videos []domian.Video) error {
	r.topN.Store(videos)
	r.ddl.Store(time.Now().Add(r.expiration))
	return nil
}

func (r *RankingLocalCache) Get(ctx context.Context) ([]domian.Video, error) {
	ddl := r.ddl.Load()
	videos := r.topN.Load()
	if len(videos) == 0 || ddl.Before(time.Now()) {
		return nil, errors.New("local cache failure")
	}
	return videos, nil
}
func (r *RankingLocalCache) ForceGet(ctx context.Context) ([]domian.Video, error) {
	return r.topN.Load(), nil
}
