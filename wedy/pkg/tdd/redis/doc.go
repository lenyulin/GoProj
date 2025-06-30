package redis

//
//import (
//	"context"
//	"github.com/alicebob/miniredis/v2"
//	"github.com/redis/go-redis/v9"
//	"strings"
//	"testing"
//	"time"
//)
//
//// miniredis是一个纯go实现的用于单元测试的redis server。
//// 它是一个简单易用的、基于内存的redis替代品，它具有真正的TCP接口，你可以把它当成是redis版本的net/http/httptest。
//
//const (
//	KeyValidWebsite = "app:valid:website:list"
//)
//
//func DoSomethingWithRedis(rdb *redis.Client, key string) bool {
//	// 这里可以是对redis操作的一些逻辑
//	ctx := context.TODO()
//	if !rdb.SIsMember(ctx, KeyValidWebsite, key).Val() {
//		return false
//	}
//	val, err := rdb.Get(ctx, key).Result()
//	if err != nil {
//		return false
//	}
//	if !strings.HasPrefix(val, "https://") {
//		val = "https://" + val
//	}
//	// 设置 blog key 五秒过期
//	if err := rdb.Set(ctx, "blog", val, 5*time.Second).Err(); err != nil {
//		return false
//	}
//	return true
//}
//
//// 下面的代码是我使用miniredis库为DoSomethingWithRedis函数编写的单元测试代码，
//// 其中miniredis不仅支持mock常用的Redis操作，还提供了很多实用的帮助函数，
//// 例如检查key的值是否与预期相等的s.CheckGet()和帮助检查key过期时间的s.FastForward()。
//func TestDoSomethingWithRedis(t *testing.T) {
//	// mock一个redis server
//	s, err := miniredis.Run()
//	if err != nil {
//		panic(err)
//	}
//	defer s.Close()
//
//	// 准备数据
//	s.Set("q1mi", "liwenzhou.com")
//	s.SAdd(KeyValidWebsite, "q1mi")
//
//	// 连接mock的redis server
//	rdb := redis.NewClient(&redis.Options{
//		Addr: s.Addr(), // mock redis server的地址
//	})
//
//	// 调用函数
//	ok := DoSomethingWithRedis(rdb, "q1mi")
//	if !ok {
//		t.Fatal()
//	}
//
//	// 可以手动检查redis中的值是否复合预期
//	if got, err := s.Get("blog"); err != nil || got != "https://liwenzhou.com" {
//		t.Fatalf("'blog' has the wrong value")
//	}
//	// 也可以使用帮助工具检查
//	s.CheckGet(t, "blog", "https://liwenzhou.com")
//
//	// 过期检查
//	s.FastForward(5 * time.Second) // 快进5秒
//	if s.Exists("blog") {
//		t.Fatal("'blog' should not have existed anymore")
//	}
//}
