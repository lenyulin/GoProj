package cache

import (
	"GoProj/wedy/interactive/domain"
	"context"
	_ "embed"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
)

var (
	//go:embed lua/incr_cnt.lua
	luaIncrCnt string
)

const fieldReadCnt = "read_cnt"

type InteractiveCache interface {
	IncrReadCntIfPresent(ctx context.Context, biz string, bizId int64) error
	Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error)
}
type InteractiveRedisCache struct {
	client redis.Cmdable
}

func NewInteractiveRedisCache(client redis.Cmdable) InteractiveCache {
	return &InteractiveRedisCache{client: client}
}
func (i *InteractiveRedisCache) IncrReadCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	return i.client.Eval(ctx, luaIncrCnt, []string{i.generateKey(biz, bizId)}, fieldReadCnt, 1).Err()
}
func (i *InteractiveRedisCache) Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error) {
	key := i.generateKey(biz, bizId)
	res, err := i.client.HGetAll(ctx, key).Result()
	if err != nil {
		return domain.Interactive{}, err
	}
	if len(res) == 0 {
		return domain.Interactive{}, ErrKeyNotExist
	}
	var v domain.Interactive
	v.BizId = bizId
	v.ReadCnt, err = strconv.ParseInt(res[fieldReadCnt], 10, 64)
	return v, err
}

func (i *InteractiveRedisCache) generateKey(biz string, bizId int64) string {
	return fmt.Sprintf("interative:video:%s:%d", biz, bizId)
}
