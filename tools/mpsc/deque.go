package mpsc

import (
	"sync/atomic"
	"unsafe"
)

type node struct {
	next *node
	val  interface{}
}

type Queue struct {
	head, tail *node
}

func New() *Queue {
	q := &Queue{}
	stub := &node{}
	q.head = stub
	q.tail = stub
	return q
}

// Push adds x to the back of the queue.
//
// Push can be safely called from multiple goroutines
func (this *Queue) Push(x interface{}) {
	n := new(node)
	n.val = x
	// current producer acquires head node
	prev := (*node)(atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&this.head)), unsafe.Pointer(n)))
	// release node to consumer
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&prev.next)), unsafe.Pointer(n))
}

// Pop removes the item from the front of the queue or nil if the queue is empty
//
// Pop must be called from a single, consumer goroutine
func (this *Queue) Pop() interface{} {
	tail := this.tail
	next := (*node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&tail.next))))
	if next != nil {
		this.tail = next
		v := next.val
		next.val = nil
		return v
	}
	return nil

}

// Empty returns true if the queue is empty
//
// Empty must be called from a single, consumer goroutine
func (this *Queue) Empty() bool {
	tail := this.tail
	next := (*node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&tail.next))))
	return next == nil
}
