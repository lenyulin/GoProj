package comment

import (
	"GoProj/wedy/comment/domain"
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func initCache(cache *localCommentCache) {
	com := domain.Comment{
		Id:      123453,
		VId:     12345,
		Content: "this is the first comment content",
		User: domain.User{
			Id:        123123,
			Name:      "123123",
			AvatarURL: "",
		}}
	err := cache.SetAndUpdate(context.Background(), strconv.Itoa(int(com.VId)), com)
	if err != nil {
		panic(err)
	}
}

func TestLocalCommentCache_SetAndUpdate(t *testing.T) {
	cache := &localCommentCache{}
	for i := 0; i < 15; i++ {
		initCache(cache)
	}
	com := domain.Comment{
		Id:      12345,
		VId:     12345,
		Content: "this is the 2rd comment content",
		User: domain.User{
			Id:        12345,
			Name:      "12345",
			AvatarURL: "",
		}}
	err := cache.SetAndUpdate(context.Background(), strconv.Itoa(int(com.VId)), com)
	if err != nil {
		fmt.Println(err)
	}
}
func TestLocalCommentCache_Get(t *testing.T) {
	cache := NewLocalCommentCache()
	time := time.Now().UnixMilli()
	coms := []domain.Comment{
		{Id: 12345,
			VId:        2222,
			Content:    "this is the first comment content",
			Dislike:    123,
			Like:       125,
			Ctime:      time,
			ReplyCount: 10,
			User: domain.User{
				Id:        1234567,
				Name:      "1234567",
				AvatarURL: "",
			}},
		{Id: 123456,
			VId:        2222,
			Dislike:    1232,
			Like:       1253,
			Ctime:      time,
			ReplyCount: 101,
			Content:    "this is the 2rd comment content",
			User: domain.User{
				Id:        12345678,
				Name:      "12345678",
				AvatarURL: "",
			}},
		{Id: 1234567,
			VId:        222333,
			Dislike:    123,
			Like:       123,
			Ctime:      time,
			ReplyCount: 123,
			Content:    "this is the 3th comment content",
			User: domain.User{
				Id:        123456789,
				Name:      "123456789",
				AvatarURL: "",
			}},
		{Id: 12345678,
			VId:        222333,
			Dislike:    123,
			Like:       1235,
			Ctime:      time,
			ReplyCount: 123,
			Content:    "this is the 4th comment content",
			User: domain.User{
				Id:        123456789,
				Name:      "123456789",
				AvatarURL: "",
			}},
	}
	for _, com := range coms {
		cache.SetAndUpdate(context.Background(), strconv.Itoa(int(com.VId)), com)
	}
	res, err := cache.Get(context.Background(), "222333", 0)
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range res {
		fmt.Println(v)
	}
	fmt.Println("try to like video")
	err = cache.IncrLikeCnt(context.Background(), 12345678, 222333, 123456789)
	if err != nil {
		fmt.Println(err)
	}
	res, err = cache.Get(context.Background(), "222333", 0)
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range res {
		fmt.Println(v)
	}
}
