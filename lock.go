package luchen

import (
	"hash/crc32"
	"sync"
)

// segmentLock 分段锁，用于减小锁粒度，避免全局锁的性能瓶颈
// 原理是将一个大锁拆分成多个小锁，不同的资源使用不同的锁，从而减少锁竞争
type segmentLock struct {
	size  int           // 锁的分段数量
	locks []*sync.Mutex // 锁数组，每个元素代表一个分段锁
}

// newSegmentLock 创建一个新的分段锁
// size: 分段数量，建议设置为 2 的幂次方，可以优化取模运算
func newSegmentLock(size int) *segmentLock {
	locks := make([]*sync.Mutex, size)
	for i := 0; i < size; i++ {
		locks[i] = &sync.Mutex{}
	}
	return &segmentLock{
		size:  size,
		locks: locks,
	}
}

// getLock 根据资源标识获取对应的分段锁
// source: 资源标识字符串
// 返回值: 该资源对应的互斥锁
func (s *segmentLock) getLock(source string) *sync.Mutex {
	hashCode := hash(source)
	idx := int(hashCode) % s.size
	return s.locks[idx]
}

// hash 计算字符串的哈希值
// 使用 crc32 算法计算字符串的哈希值，用于确定使用哪个分段锁
func hash(s string) uint32 {
	return crc32.ChecksumIEEE([]byte(s))
}
