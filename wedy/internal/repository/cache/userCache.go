package cache

import (
	"GoProj/wedy/internal/domian"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
	"time"
)

type RedisUserCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

type UserCache interface {
	Get(ctx context.Context, phone string) (domian.User, error)
	Set(ctx context.Context, du domian.User) error
}

var ErrKeyNotExists = redis.Nil

func NewUserCache(cmd redis.Cmdable) UserCache {
	return &RedisUserCache{
		cmd:        cmd,
		expiration: time.Minute * 15,
	}
}

func (uc *RedisUserCache) Get(ctx context.Context, phone string) (domian.User, error) {
	key := uc.key(phone)
	data, err := uc.cmd.Get(ctx, key).Result()
	if err != nil {
		return domian.User{}, err
	}
	var u domian.User
	err = json.Unmarshal([]byte(data), &u)
	return u, err
}
func (uc *RedisUserCache) key(phone string) string {
	return fmt.Sprintf("user:info:%s", phone)
}

func (uc *RedisUserCache) Set(ctx context.Context, du domian.User) error {
	key := uc.key(du.Phone)
	data, err := json.Marshal(du)
	if err != nil {
		return err
	}
	return uc.cmd.Set(ctx, key, string(data), uc.expiration).Err()
}
