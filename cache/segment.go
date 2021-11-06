package cache

import (
	"errors"
	"fmt"
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
	rb            reflect.Value
	segId         int
	_             uint32
	missCount     int64
	hitCount      int64
	entryCount    int64
	totalCount    int64      // number of entries in ring buffer, including deleted entries.
	totalTime     int64      // used to calculate least recent used entry.
	timer         Timer      // Timer giving current time
	totalEvacuate int64      // used for debug
	totalExpired  int64      // used for debug
	overwrites    int64      // used for debug
	touched       int64      // used for debug
	vacuumLen     int64      // up to vacuumLen, new data can be written without overwriting old data.
	slotLens      [256]int32 // The actual length for every slot.
	slotCap       int32      // max number of entry pointers a slot can hold.
	slotsData     []entryPtr // shared by all 256 slots
}

func newSegment(bufSize int, segId int, timer Timer) (seg segment) {
	seg.rb = NewRingBuf(bufSize, 0)
	seg.segId = segId
	seg.timer = timer
	seg.vacuumLen = int64(bufSize)
	seg.slotCap = 1
	seg.slotsData = make([]entryPtr, 256*seg.slotCap)
	return
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
		fmt.Println(matchedPtr)
	} else {
		hdr.slotId = slotId
		hdr.hash16 = hash16
		hdr.keyLen = uint16(len(key))
		hdr.accessTime = now
		hdr.expireAt = expireAt
		hdr.valLen = uint32(valueOf.Len())
		hdr.valCap = uint32(valueOf.Len())
		if hdr.valCap == 0 { // avoid infinite loop when increasing capacity.
			hdr.valCap = 1
		}
	}
	entryLen := ENTRY_HDR_SIZE + int64(len(key)) + int64(hdr.valCap)
	slotModified := this.evacuate(entryLen, slotId, now)
	if slotModified {
		slot = this.getSlot(slotId)
		idx, match = this.lookup(slot, hash16, key)
	}
	newOff := this.rb.Len()
	this.insertEntryPtr(slotId, hash16, int64(newOff), idx, hdr.keyLen)
	atomic.AddInt64(&this.totalTime, int64(now))
	atomic.AddInt64(&this.totalCount, 1)
	this.vacuumLen -= entryLen

	return
}

func (this *segment) evacuate(entryLen int64, slotId uint8, now uint32) (slotModified bool) {
	var oldHdrBuf [ENTRY_HDR_SIZE]byte
	consecutiveEvacuate := 0
	for this.vacuumLen < entryLen {
		oldHdr := (*entryHdr)(unsafe.Pointer(&oldHdrBuf[0]))
		oldEntryLen := ENTRY_HDR_SIZE + int64(oldHdr.keyLen) + int64(oldHdr.valCap)
		if oldHdr.deleted {
			consecutiveEvacuate = 0
			atomic.AddInt64(&this.totalTime, -int64(oldHdr.accessTime))
			atomic.AddInt64(&this.totalCount, -1)
			this.vacuumLen += oldEntryLen
			continue
		}
		expired := oldHdr.expireAt != 0 && oldHdr.expireAt < now
		leastRecentUsed := int64(oldHdr.accessTime)*atomic.LoadInt64(&this.totalCount) <= atomic.LoadInt64(&this.totalTime)
		if expired || leastRecentUsed || consecutiveEvacuate > 5 {
			if oldHdr.slotId == slotId {
				slotModified = true
			}
			consecutiveEvacuate = 0
			atomic.AddInt64(&this.totalTime, -int64(oldHdr.accessTime))
			atomic.AddInt64(&this.totalCount, -1)
			this.vacuumLen += oldEntryLen
			if expired {
				atomic.AddInt64(&this.totalExpired, 1)
			} else {
				atomic.AddInt64(&this.totalEvacuate, 1)
			}
		} else {
			// evacuate an old entry that has been accessed recently for better cache hit rate.
			consecutiveEvacuate++
			atomic.AddInt64(&this.totalEvacuate, 1)
		}
	}
	return
}

func (this *segment) updateEntryPtr(slotId uint8, hash16 uint16, oldOff, newOff int64) {
	slot := this.getSlot(slotId)
	idx, match := this.lookupByOff(slot, hash16, oldOff)
	if !match {
		return
	}
	ptr := &slot[idx]
	ptr.offset = newOff
}

func (seg *segment) lookupByOff(slot []entryPtr, hash16 uint16, offset int64) (idx int, match bool) {
	idx = entryPtrIdx(slot, hash16)
	for idx < len(slot) {
		ptr := &slot[idx]
		if ptr.hash16 != hash16 {
			break
		}
		match = ptr.offset == offset
		if match {
			return
		}
		idx++
	}
	return
}
func (this *segment) expand() {
	newSlotData := make([]entryPtr, this.slotCap*2*256)
	for i := 0; i < 256; i++ {
		off := int32(i) * this.slotCap
		copy(newSlotData[off*2:], this.slotsData[off:off+this.slotLens[i]])
	}
	this.slotCap *= 2
	this.slotsData = newSlotData
}

func (this *segment) insertEntryPtr(slotId uint8, hash16 uint16, offset int64, idx int, keyLen uint16) {
	if this.slotLens[slotId] == this.slotCap {
		this.expand()
	}
	this.slotLens[slotId]++
	atomic.AddInt64(&this.entryCount, 1)
	slot := this.getSlot(slotId)
	copy(slot[idx+1:], slot[idx:])
	slot[idx].offset = offset
	slot[idx].hash16 = hash16
	slot[idx].keyLen = keyLen
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
