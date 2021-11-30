package cache

import (
	"sync"

	"github.com/cespare/xxhash"
)

const (
	// segmentCount represents the number of segments within a freecache instance.
	segmentCount = 256
	// segmentAndOpVal is bitwise AND applied to the hashVal to find the segment id.
	segmentAndOpVal = 255
	minBufSize      = 512 * 1024
)

// Cache is a freecache instance.
type Cache struct {
	locks    [segmentCount]sync.Mutex
	segments [segmentCount]segment
}

func hashFunc(data []byte) uint64 {
	return xxhash.Sum64(data)
}

func NewCache(size int) (cache *Cache) {
	return NewCacheCustomTimer(size, defaultTimer{})
}

func NewCacheCustomTimer(size int, timer Timer) (cache *Cache) {
	if size < minBufSize {
		size = minBufSize
	}
	if timer == nil {
		timer = defaultTimer{}
	}
	cache = new(Cache)
	for i := 0; i < segmentCount; i++ {
		//cache.segments[i] = newSegment(size/segmentCount, timer)
	}
	return
}
