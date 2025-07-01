package comment

import (
	"GoProj/wedy/comment/domain"
	"context"
	"errors"
	"fmt"
	avl "github.com/emirpasic/gods/trees/avltree"
	"github.com/redis/go-redis/v9"
	"strconv"
	"sync"
	"time"
)

var (
	Expriation = time.Minute * 10
)

func NewLocalCommentCache(redis redis.Cmdable) Cache {
	return &localCommentCache{
		batchSize: 20,
		getScore: func(duration int64, like int64, dislike int64, reply int64) int64 {
			return (time.Now().UnixMilli() - duration) + like - dislike + reply
		},
		vComments: make(map[string]*CachedComment),
		redisCache: redisCommentCache{
			client:       redis,
			expatriation: time.Minute * 60,
			batchSize:    20,
		},
	}
}

var (
	ErrVideoRecordNotFound = errors.New("video record not found")
)

type CachedComment struct {
	ddl      time.Time
	comments map[int64]*domain.Comment
	count    int64
	score    *avl.Tree
}

type localCommentCache struct {
	mu         sync.RWMutex
	vComments  map[string]*CachedComment
	cachedVid  []int64
	batchSize  int64
	getScore   func(duration int64, like int64, dislike int64, reply int64) int64
	redisCache redisCommentCache
}

func (l *localCommentCache) SetCommentLikeTop(ctx context.Context) error {
	vCmms := make(map[string]*CachedComment, len(l.cachedVid))
	for _, vid := range l.cachedVid {
		res, err := l.redisCache.Get(ctx, strconv.FormatInt(vid, 10), 0)
		if err == nil {
			cmms := make(map[int64]*domain.Comment)
			for _, v := range res {
				cmms[v.Id] = &v
			}
			lCache := &CachedComment{
				ddl:      time.Now().Add(Expriation),
				comments: cmms,
				count:    int64(len(res)),
			}
			vCmms[strconv.FormatInt(vid, 10)] = lCache
		}
	}
	return nil
}

func (l *localCommentCache) IncrLikeCnt(ctx context.Context, id int64, vid int64, uid int64) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	res, ok := l.vComments[strconv.FormatInt(vid, 10)]
	if !ok {
		return ErrVideoRecordNotFound
	}
	if re, loaded := res.comments[id]; loaded {
		re.Like += 1
		res.comments[id] = re
		l.vComments[strconv.FormatInt(vid, 10)] = res
		return nil
	}
	return ErrCommentRecordNotFound
}

func (l *localCommentCache) SetAndUpdate(ctx context.Context, key string, comment domain.Comment) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	res, ok := l.vComments[key]
	if !ok {
		//未获取记录，记录为空，首次访问，写入首个评论
		newCmm := &CachedComment{
			ddl:   time.Now().Add(Expriation),
			count: 1,
			score: avl.NewWithIntComparator(),
		}
		newCmm.comments = make(map[int64]*domain.Comment)
		newCmm.comments[comment.Id] = &comment
		s := int(l.getScore(comment.Ctime, comment.Like, comment.Dislike, comment.ReplyCount))
		newCmm.score.Put(s, comment.Id)
		l.cachedVid = append(l.cachedVid, comment.VId)
		l.vComments[key] = newCmm
		return nil
	}
	// 检查评论是否已存在，异常，理论不应该存在
	if _, loaded := res.comments[comment.Id]; loaded {
		fmt.Printf("comment id exist, id:%d\n", comment.Id)
		return nil
	}
	// 评论不存在，添加新评论
	s := int(l.getScore(comment.Ctime, comment.Like, comment.Dislike, comment.ReplyCount))
	_, ok = res.score.Get(s)
	if !ok {
		res.count += 1
		res.score.Put(s, comment.Id)
		res.ddl = time.Now().Add(Expriation)
		res.comments[comment.Id] = &comment
		l.vComments[key] = res
		return nil
	}
	return nil
}
func (l *localCommentCache) DeleteExpiredComments() {

}

//func (l *localCommentCache) BatchSetAndUpdate(ctx context.Context, key string, comment []domain.Comment) error {
//	res, _ := l.comments.LoadOrStore(key, &atomic.Value{})
//	val := res.(*atomic.Value)
//	oldCmm, ok := val.Load().(*CachedComment)
//	if !ok {
//		//未获取记录，记录为空，首次访问，写入首个评论
//		newCmm := &CachedComment{
//			ddl: time.Now().Add(l.expiration),
//			comments: []domain.Comment{
//				comment,
//			},
//		}
//		val.Store(newCmm)
//		return nil
//	}
//	oldCmm.mu.Lock()
//	defer oldCmm.mu.Unlock()
//	newComments := make([]domain.Comment, len(oldCmm.comments))
//	copy(newComments, oldCmm.comments)
//	newComments = append(newComments, comment)
//	newCachedComment := &CachedComment{
//		ddl:      time.Now().Add(l.expiration),
//		comments: newComments,
//	}
//	newCachedComment.count = oldCmm.count + 1
//	val.Store(newCachedComment)
//	return nil
//}

//	func (l *localCommentCache) listComments(key string) {
//		res, ok := l.comments.Load(key)
//		if !ok {
//			fmt.Println("no comment found")
//		}
//		val := res.(*atomic.Value)
//		av, ok := val.Load().(*CachedComment)
//		if ok {
//			fmt.Printf("%p\n", &av)
//		}
//	}

func (l *localCommentCache) Get(ctx context.Context, key string, offset int64) ([]domain.Comment, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	val, ok := l.vComments[key]
	if !ok {
		return nil, ErrCommentRecordNotFound
	}
	if val.ddl.Before(time.Now()) {
		//缓存过期处理
		return nil, ErrCommentCacheExpired
	}
	if val.count-offset <= 0 {
		return nil, ErrCommentRecordNotFound
	}
	index := val.score.Values()[offset:l.getMaxOffset(offset, l.batchSize, val.count)]
	var res = make([]domain.Comment, 0)
	for _, v := range index {
		re, _ := val.comments[v.(int64)]
		res = append(res, *re)
	}
	return res, nil
}
func (l *localCommentCache) getMaxOffset(offset int64, batch int64, commentCount int64) int64 {
	if commentCount-offset < batch && commentCount-offset > 0 {
		return commentCount
	}
	return offset + batch
}
