package cache

import (
	"errors"
	"reflect"
	"sync/atomic"
	"unsafe"
)

const HASH_ENTRY_SIZE = 16
const ENTRY_HDR_SIZE = 24

var ErrLargeKey = errors.New("The key is larger than 65535")
var ErrLargeEntry = errors.New("The entry size is larger than 1/1024 of cache size")
var ErrNotFound = errors.New("Entry not found")

// entry pointer struct points to an entry in ring buffer
type entryPtr struct {
	offset   int64  // entry offset in ring buffer
	hash16   uint16 // entries are ordered by hash16 in a slot.
	keyLen   uint16 // used to compare a key
	reserved uint32
}

// entry header struct in ring buffer, followed by key and value.
type entryHdr struct {
	accessTime uint32
	expireAt   uint32
	keyLen     uint16
	hash16     uint16
	valLen     uint32
	valCap     uint32
	deleted    bool
	slotId     uint8
	reserved   uint16
}

// a segment contains 256 slots, a slot is an array of entry pointers ordered by hash16 value
// the entry can be looked up by hash value of the key.
type segment struct {
	rb reflect.Value

	totalTime int64 // used to calculate least recent used entry.
	timer     Timer // Timer giving current time

	slotLens  [256]int32 // The actual length for every slot.
	slotCap   int32      // max number of entry pointers a slot can hold.
	slotsData []entryPtr // shared by all 256 slots
}

func (this segment) set(key []byte, value interface{}, hashVal int64, expireSeconds int) (err error) {
	valueOf := reflect.ValueOf(value)
	if len(key) > 65535 {
		return ErrLargeKey
	}
	now := this.timer.Now()
	var expireAt uint32 = 0
	if expireSeconds > 0 {
		expireAt = now + uint32(expireSeconds)
	}
	slotId := uint8(hashVal >> 8)
	hash16 := uint16(hashVal >> 16)
	slot := this.getSlot(slotId)
	idx, match := this.lookup(slot, hash16, key)

	var hdrBuf [ENTRY_HDR_SIZE]byte
	hdr := (*entryHdr)(unsafe.Pointer(&hdrBuf[0]))
	if match {
		matchedPtr := &slot[idx]
		hdr.slotId = slotId
		hdr.hash16 = hash16
		hdr.keyLen = uint16(len(key))
		originAccessTime := hdr.accessTime
		hdr.accessTime = now
		hdr.expireAt = expireAt
		hdr.valLen = uint32(valueOf.Len())
		if hdr.valCap >= hdr.valLen {
			atomic.AddInt64(&this.totalTime, int64(hdr.accessTime)-int64(originAccessTime))
		}
	}
	return
}
func entryPtrIdx(slot []entryPtr, hash16 uint16) (idx int) {
	high := len(slot)
	for idx < high {
		mid := (idx + high) >> 1
		oldEntry := &slot[mid]
		if oldEntry.hash16 < hash16 {
			idx = mid + 1
		} else {
			high = mid
		}
	}
	return
}
func (this *segment) lookup(slot []entryPtr, hash16 uint16, key []byte) (idx int, match bool) {
	idx = entryPtrIdx(slot, hash16)
	for idx < len(slot) {
		ptr := &slot[idx]
		if ptr.hash16 != hash16 {
			break
		}
		match = int(ptr.keyLen) == len(key) && this.rb.Bool()
		if match {
			return
		}
		idx++
	}
	return
}

func (this *segment) getSlot(slotId uint8) []entryPtr {
	slotOff := int32(slotId) * this.slotCap
	return this.slotsData[slotOff : slotOff+this.slotLens[slotId] : slotOff+this.slotCap]
}
