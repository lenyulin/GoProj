package bigcachex

import (
	"context"
	"errors"
	"fmt"
	"github.com/allegro/bigcache/v3"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

type bigcachex struct {
	cache  *bigcache.BigCache
	redis  redis.Cmdable
	rwLock sync.RWMutex
	biz    string
	bizId  string
}
type redisCache struct {
	Data    []any
	version uint64
}

func NewBigCachex(redis redis.Cmdable) BigCachex {
	cache, err := bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Minute))
	if err != nil {
		panic(err)
	}
	return &bigcachex{
		cache: cache,
		redis: redis,
	}
}

func (b *bigcachex) Serve() error {
	return nil
}

func (b *bigcachex) Shutdown() error {
	b.rwLock.Lock()
	defer b.rwLock.Unlock()
	err := b.cache.Close()
	return err
}

func (b *bigcachex) Get(key string) ([]byte, error) {
	entry, err := b.cache.Get(key)
	if err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
			res, err := b.redis.Get(ctx, b.generateKey(key)).Result()
			if err != nil {
				return nil, err
			}
			b.rwLock.Lock()
			defer b.rwLock.Unlock()
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

func (b *bigcachex) Set(key string, value []byte) error {
	err := b.cache.Set(key, value)
	if err != nil {
		return err
	}
	return nil
}

func (b *bigcachex) Delete(key string) error {
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
