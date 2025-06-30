package repository

import (
	"GoProj/wedy/internal/domian"
	"GoProj/wedy/internal/repository/cache"
	"context"
)

type RankingRepository interface {
	ReplaceRedisTopN(ctx context.Context, videos []domian.Video) error
	GetRedisTopN(ctx context.Context) ([]domian.Video, error)
	GetTopN(ctx context.Context) ([]domian.Video, error)
}

func (r *CacheRankingRepository) GetRedisTopN(ctx context.Context) ([]domian.Video, error) {
	return r.redisCache.Get(ctx)
}
func (r *CacheRankingRepository) GetTopN(ctx context.Context) ([]domian.Video, error) {
	res, err := r.redisCache.Get(ctx)
	if err == nil {
		return res, nil
	}
	res, err = r.redisCache.Get(ctx)
	if err != nil {
		return r.localCache.ForceGet(ctx)
	}
	_ = r.localCache.Set(ctx, res)
	return res, nil
}

type CacheRankingRepository struct {
	cache      cache.RankingCache
	redisCache *cache.RedisRankingCache
	localCache *cache.RankingLocalCache
}

func NewCacheRankingRepository(cache cache.RankingCache) RankingRepository {
	return &CacheRankingRepository{cache: cache}
}

func (r *CacheRankingRepository) ReplaceRedisTopN(ctx context.Context, videos []domian.Video) error {
	_ = r.localCache.Set(ctx, videos)
	return r.cache.Set(ctx, videos)
}
