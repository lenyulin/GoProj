package cache

import (
	"GoProj/wedy/internal/domian"
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"time"
)

type RankingCache interface {
	Set(ctx context.Context, videos []domian.Video) error
	Get(ctx context.Context) ([]domian.Video, error)
}
type RedisRankingCache struct {
	client       *redis.Client
	key          string
	expatriation time.Duration
}

func NewRedisRankingCache(client *redis.Client) RankingCache {
	return &RedisRankingCache{client: client, expatriation: time.Second * 30, key: "ranking:topN"}
}
func (c *RedisRankingCache) Get(ctx context.Context) ([]domian.Video, error) {
	val, errl := c.client.Get(ctx, c.key).Result()
	if errl != nil {
		return nil, errl
	}
	var res []domian.Video
	err := json.Unmarshal([]byte(val), &res)
	return res, err
}
func (c *RedisRankingCache) Set(ctx context.Context, videos []domian.Video) error {
	val, err := json.Marshal(&videos)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, c.key, val, c.expatriation).Err()
}
