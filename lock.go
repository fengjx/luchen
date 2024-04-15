package luchen

import (
	"hash/crc32"
	"sync"
)

// segmentLock 分段锁，减小锁粒度
type segmentLock struct {
	size    int
	lockMap map[int]*sync.Mutex
}

func newSegmentLock(size int) *segmentLock {
	lockMap := make(map[int]*sync.Mutex, size)
	for i := 0; i < size; i++ {
		lockMap[i] = &sync.Mutex{}
	}
	return &segmentLock{
		size:    size,
		lockMap: lockMap,
	}
}

func (s *segmentLock) getLock(source string) *sync.Mutex {
	hashCode := hash(source)
	idx := int(hashCode) % s.size
	return s.lockMap[idx]
}

func hash(s string) uint32 {
	return crc32.ChecksumIEEE([]byte(s))
}
