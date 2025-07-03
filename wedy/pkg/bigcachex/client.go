package bigcachex

import (
	"GoProj/wedy/pkg/bigcachex/proto"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/allegro/bigcache/v3"
	"github.com/redis/go-redis/v9"
	proto2 "google.golang.org/protobuf/proto"
	"time"
)

type bigcachex struct {
	cache *bigcache.BigCache
	redis redis.Cmdable
	biz   string
	bizId string
}
type redisCache struct {
	Data    []any
	version uint64
}

func NewBigCachex(redis redis.Cmdable, biz string, bizId string) HybridCache {
	cache, err := bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Minute))
	if err != nil {
		panic(err)
	}
	return &bigcachex{
		cache: cache,
		redis: redis,
		biz:   biz,
		bizId: bizId,
	}
}

func (b *bigcachex) Serve() error {
	return nil
}

func (b *bigcachex) Shutdown() error {
	err := b.cache.Close()
	return err
}

func (b *bigcachex) Get(ctx context.Context, key string) ([]byte, error) {
	entry, err := b.cache.Get(key)
	if err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			res, err := b.redis.Get(ctx, b.generateKey(key)).Result()
			if err != nil {
				return nil, err
			}
			err = b.cache.Set(key, []byte(res))
			if err != nil {
				return nil, err
			}
			return []byte(res), nil
		}
		return nil, err
	}
	return entry, nil
}

//go:embed update_cache.lua
var update_cache string

func (b *bigcachex) Set(ctx context.Context, key string, value []byte) error {
	err := b.cache.Set(key, value)
	if err != nil {
		return err
	}
	go func() {
		var act proto.SeckillActivity
		err := proto2.Unmarshal(value, &act)
		if err != nil {
			err = b.redis.Eval(ctx, update_cache, []string{b.bizId}, act.Version, value).Err()
		}
	}()
	return nil
}

func (b *bigcachex) Delete(ctx context.Context, key string) error {
	err := b.cache.Delete(key)
	if err != nil {
		return err
	}
	return nil
}
func (b *bigcachex) close() error {
	err := b.cache.Close()
	if err != nil {
		return err
	}
	return nil
}
func (b *bigcachex) generateKey(key string) string {
	return fmt.Sprintf("biz:%s:bizId:%s:key%s", b.biz, b.bizId, key)
}
