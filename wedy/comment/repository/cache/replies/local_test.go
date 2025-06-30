package replies

import (
	"GoProj/wedy/comment/domain"
	"context"
	"fmt"
	"strconv"
	"testing"
)

func initCache(cache *localReplyCache) {
	com := domain.Reply{
		Id:      123453,
		ReplyTo: 12345,
		Content: "this is the first comment content",
		ReplyToUser: domain.ReplyToUser{
			Id:   12345,
			Name: "12345",
		}}
	err := cache.SetAndUpdate(context.Background(), strconv.Itoa(int(com.ReplyTo)), com)
	if err != nil {
		panic(err)
	}
}

func TestLocalCommentCache_SetAndUpdate(t *testing.T) {
	cache := NewLocalReplyCache()
	for i := 0; i < 15; i++ {
		com := domain.Reply{
			Id:      123453,
			ReplyTo: 12345,
			Content: "this is the first comment content",
			ReplyToUser: domain.ReplyToUser{
				Id:   12345,
				Name: "12345",
			}}
		cache.SetAndUpdate(context.Background(), strconv.Itoa(int(com.ReplyTo)), com)
	}
	com := domain.Reply{
		Id:      123456,
		ReplyTo: 12345,
		Content: "this is the 2rd comment content",
		ReplyToUser: domain.ReplyToUser{
			Id:   12345,
			Name: "12345",
		}}
	err := cache.SetAndUpdate(context.Background(), strconv.Itoa(int(com.ReplyTo)), com)
	if err != nil {
		fmt.Println(err)
	}
}
func TestLocalCommentCache_Get(t *testing.T) {
	cache := NewLocalReplyCache()
	for i := 0; i < 10; i++ {
		com := domain.Reply{
			Id:      123453,
			ReplyTo: 12345,
			Content: "this is the first comment content",
			ReplyToUser: domain.ReplyToUser{
				Id:   12345,
				Name: "12345",
			}}
		cache.SetAndUpdate(context.Background(), strconv.Itoa(int(com.ReplyTo)), com)
	}
	res, err := cache.Get(context.Background(), "12345", 10)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(len(res))
	fmt.Println(res)
}
