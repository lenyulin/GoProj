package cache

import (
	"GoProj/wedy/internal/domian"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type VideoCache interface {
	GetFirstPage(ctx context.Context, id int64) ([]domian.Video, error)
	SetFirstPage(ctx context.Context, id int64, res []domian.Video) error
	Get(ctx context.Context, id int64) (domian.Video, error)
	Set(ctx context.Context, video domian.Video) error
}
type VideoRedisCache struct {
	redis redis.Cmdable
}

func (v *VideoRedisCache) Set(ctx context.Context, video domian.Video) error {
	key := v.videoAuthorKey(video.VUid)
	res, err := json.Marshal(video)
	if err != nil {
		return err
	}
	err = v.redis.Set(ctx, key, string(res), 0).Err()
	if err != nil {
		return err
	}
	return nil
}
func (v *VideoRedisCache) Get(ctx context.Context, id int64) (domian.Video, error) {
	key := v.videoAuthorKey(id)
	res, err := v.redis.Get(ctx, key).Bytes()
	if err != nil {
		return domian.Video{}, err
	}
	var videos []domian.Video
	err = json.Unmarshal(res, &videos)
	if err != nil {
		return domian.Video{}, err
	}
	return videos[0], err
}
func NewVideoRedisCache(redis redis.Cmdable) VideoCache {
	return &VideoRedisCache{
		redis: redis,
	}
}
func (v *VideoRedisCache) GetFirstPage(ctx context.Context, id int64) ([]domian.Video, error) {
	key := v.key(id)
	res, err := v.redis.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	var videos []domian.Video
	err = json.Unmarshal(res, &videos)
	if err != nil {
		return nil, err
	}
	return videos, err
}

func (v *VideoRedisCache) SetFirstPage(ctx context.Context, id int64, videos []domian.Video) error {
	key := v.key(id)
	bytes, err := json.Marshal(videos)
	if err != nil {
		return err
	}
	err = v.redis.Set(ctx, key, bytes, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (v *VideoRedisCache) key(id int64) string {
	return fmt.Sprintf("videos:firstpage:%d", id)
}
func (v *VideoRedisCache) videoAuthorKey(id int64) string {
	return fmt.Sprintf("videos:authorid:%d", id)
}
