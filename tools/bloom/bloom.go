package bloom

import (
	"math"
	"sync"
)

const (
	mod7       = 1<<3 - 1
	bitPerByte = 8
)

// Filter is the sturct of BloomFilter
// false positive error rate p approximately (1 - e^(-kn/m))^k
// probability of false positives decreases as m increases, and increases as n increases.
// k is number of hash function,
// m is the size of filter, n is the number of elements inserted
type Filter struct {
	lock       *sync.RWMutex
	concurrent bool

	m     uint64 // bit array of m bits, m will be ceiling to power of 2
	n     uint64 // number of inserted elements
	log2m uint64 // log_2 of m
	k     uint64 // the number of hash function
	keys  []byte // byte array to store hash value
}

func New(size uint64, k uint64, race bool) *Filter {
	log2 := uint64(math.Ceil(math.Log2(float64(size))))

	filter := &Filter{
		m:          1 << log2,
		log2m:      log2,
		k:          k,
		keys:       make([]byte, 1<<log2),
		concurrent: race,
	}
	if filter.concurrent {
		filter.lock = &sync.RWMutex{}
	}
	return filter
}

func (f *Filter) Add(data []byte) *Filter {

	return f
}

func baseHash(data []byte) []uint64 {

	return []uint64{}
}
