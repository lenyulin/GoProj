package comment

import (
	"GoProj/wedy/comment/domain"
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"time"
)

type redisCommentCache struct {
	client       redis.Cmdable
	expatriation time.Duration
	batchSize    int64
}

func (r *redisCommentCache) IncrLikeCnt(ctx context.Context, id int64, vid int64, uid int64) error {
	//TODO implement me
	panic("implement me")
}

func (r *redisCommentCache) SetAndUpdate(ctx context.Context, key string, comment domain.Comment) error {
	val, err := json.Marshal(comment)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, val, r.expatriation).Err()
}
func (r *redisCommentCache) BatchSetAndUpdate(ctx context.Context, key string, comment domain.Comment) error {
	val, err := json.Marshal(comment)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, val, r.expatriation).Err()
}
func NewRedisCommentCache(client redis.Cmdable) Cache {
	return &redisCommentCache{
		client:       client,
		expatriation: time.Minute * 60,
		batchSize:    20,
	}
}

func (r *redisCommentCache) Get(ctx context.Context, key string, offset int64) ([]domain.Comment, error) {
	//TODO implement me
	panic("implement me")
}
