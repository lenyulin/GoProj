package domain

import (
	"math/rand"
	"sync/atomic"
	"time"
)

type Video struct {
	ID    int64        // 视频ID
	Score atomic.Int64 // 热度分数（原子操作）
}

// SkipNode 跳表节点
type SkipNode struct {
	Video *Video      // 视频指针
	next  []*SkipNode // 多层索引指针
}

// SkipList 跳表结构
type SkipList struct {
	head     *SkipNode    // 头节点
	maxLevel int          // 最大层数
	level    int          // 当前层数
	size     int          // 元素数量
	mu       atomic.Int32 // 用于写操作的轻量级锁（原子操作实现）
}

func NewSkipList(maxLevel int) *SkipList {
	head := &SkipNode{
		Video: &Video{ID: -1, Score: atomic.Int64{}},
		next:  make([]*SkipNode, maxLevel),
	}
	return &SkipList{
		head:     head,
		maxLevel: maxLevel,
		level:    1,
	}
}

// randomLevel 生成随机层数（幂次定律分布）
func (s *SkipList) randomLevel() int {
	level := 1
	for rand.Float64() < 0.5 && level < s.maxLevel {
		level++
	}
	return level
}

// Insert 插入或更新视频（按热度降序排列）
func (s *SkipList) Insert(video *Video) {
	// 使用原子操作实现轻量级锁（替代互斥锁）
	for !s.mu.CompareAndSwap(0, 1) {
		time.Sleep(time.Microsecond)
	}
	defer s.mu.Store(0)

	// 查找插入位置
	update := make([]*SkipNode, s.maxLevel)
	curr := s.head

	for i := s.level - 1; i >= 0; i-- {
		for curr.next[i] != nil && curr.next[i].Video.Score.Load() > video.Score.Load() {
			curr = curr.next[i]
		}
		update[i] = curr
	}

	// 生成随机层数
	newLevel := s.randomLevel()
	if newLevel > s.level {
		for i := s.level; i < newLevel; i++ {
			update[i] = s.head
		}
		s.level = newLevel
	}

	// 创建新节点并插入
	newNode := &SkipNode{
		Video: video,
		next:  make([]*SkipNode, newLevel),
	}

	for i := 0; i < newLevel; i++ {
		newNode.next[i] = update[i].next[i]
		update[i].next[i] = newNode
	}

	s.size++
}

// GetTopK 获取热度最高的K个视频（非阻塞读）
func (s *SkipList) GetTopK(k int) []*Video {
	result := make([]*Video, 0, k)
	curr := s.head.next[0] // 从第一层开始遍历（最高热度）

	for curr != nil && len(result) < k {
		result = append(result, curr.Video)
		curr = curr.next[0]
	}
	return result
}
