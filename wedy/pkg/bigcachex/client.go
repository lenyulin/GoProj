package bigcachex

import (
	"errors"
	"github.com/allegro/bigcache/v3"
	"github.com/withlin/canal-go/client"
	"sync"
)

type bigcachex struct {
	cache  bigcache.BigCache
	canal  client.CanalConnector
	rwLock sync.RWMutex
}

func NewBigCachex(cache bigcache.BigCache, canal client.CanalConnector) BigCachex {
	return &bigcachex{
		cache: cache,
		canal: canal,
	}
}

func (b *bigcachex) Serve() error {
	err := b.canal.Connect()
	if err != nil {
		return err
	}
}

func (b *bigcachex) Shutdown() error {
	b.rwLock.Lock()
	defer b.rwLock.Unlock()
	var errs error
	err := b.canal.DisConnection()
	if err != nil {
		errors.Join(errs, err)
	}
	err = b.cache.Close()
	if err != nil {
		errors.Join(errs, err)
	}
	return errs
}

func (b *bigcachex) Get(key string) ([]byte, error) {
	entry, err := b.cache.Get(key)
	if err != nil {
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
