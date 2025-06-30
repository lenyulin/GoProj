package startup

import (
	"GoProj/wedy/config"
	rclock "github.com/gotomicro/redis-lock"
	"github.com/redis/go-redis/v9"
)

func InitRedis() redis.Cmdable {
	return redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
}
func InitRLockClent(c redis.Cmdable) *rclock.Client {
	return rclock.NewClient(c)
}
