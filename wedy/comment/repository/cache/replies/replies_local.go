package replies

import (
	"GoProj/wedy/comment/domain"
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func NewLocalReplyCache() Cache {
	return &localReplyCache{
		expiration: time.Minute * 10,
		batchSize:  20,
	}
}

type localReplyCache struct {
	expiration time.Duration
	replies    sync.Map
	batchSize  int64
}
type CachedReply struct {
	mu      sync.RWMutex
	ddl     time.Time
	replies []domain.Reply
	count   int64
}

func (l *localReplyCache) DeleteExpiredComments() {
	panic("implement me")
}

func (l *localReplyCache) SetAndUpdate(ctx context.Context, key string, reply domain.Reply) error {
	res, _ := l.replies.LoadOrStore(key, &atomic.Value{})
	val := res.(*atomic.Value)
	oldCmm, ok := val.Load().(*CachedReply)
	if !ok {
		//未获取记录，记录为空，首次访问，写入首个评论
		newCmm := &CachedReply{
			ddl: time.Now().Add(l.expiration),
			replies: []domain.Reply{
				reply,
			},
		}
		val.Store(newCmm)
		return nil
	}
	oldCmm.mu.Lock()
	defer oldCmm.mu.Unlock()
	newReplies := make([]domain.Reply, len(oldCmm.replies))
	copy(newReplies, oldCmm.replies)
	newReplies = append(newReplies, reply)
	newCachedComment := &CachedReply{
		ddl:     time.Now().Add(l.expiration),
		replies: newReplies,
	}
	newCachedComment.count = oldCmm.count + 1
	val.Store(newCachedComment)
	return nil
}

func (l *localReplyCache) listComments(key string) {
	res, ok := l.replies.Load(key)
	if !ok {
		fmt.Println("no comment found")
	}
	val := res.(*atomic.Value)
	av, ok := val.Load().(*CachedReply)
	if ok {
		fmt.Println(av)
	}
}

func (l *localReplyCache) Get(ctx context.Context, key string, offset int64) ([]domain.Reply, error) {
	l.listComments(key)
	val, ok := l.replies.Load(key)
	if !ok {
		return nil, ErrReplyRecordNotFound
	}
	aVal := val.(*atomic.Value)
	cmm, ok := aVal.Load().(*CachedReply)
	if !ok {
		return nil, ErrLoadAtomicReply
	}
	if cmm.ddl.Before(time.Now()) {
		//缓存过期处理
		l.replies.Delete(key)
		return nil, ErrReplyCacheExpired
	}
	cmm.mu.RLock()
	defer cmm.mu.RUnlock()
	if cmm.count-offset <= 0 {
		return nil, ErrReplyRecordNotFound
	}
	return cmm.replies[offset : offset+l.batchSize-cmm.count], nil
}
