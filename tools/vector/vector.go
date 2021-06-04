package vector

import (
	"github.com/erDong01/micro-kit/core/containers"
	"log"
)

const (
	VectorBlockSize = 16
)

func assert(x bool, y string) {
	if bool(x) == false {
		log.Printf("\nFatal :{%s}", y)
	}
}

type (
	Vector struct {
		ElementCount int
		ArraySize    int
		Array        []interface{}
	}
	IVector interface {
		containers.Container
		insert(int)
		increment()
		decrement()

		Erase(int)
		PushFront(interface{})
		PushBack(interface{})
		PopFront()
		PopBack()
		Front() interface{}
		Back() interface{}
		Len() int
		Get(int) interface{}
		Swap(i, j int)
		Less(i, j int) bool
	}
)

func (this *Vector) insert(index int) {
	assert(index <= this.ElementCount, "Vector<T>::insert - out of bounds index.")
	if this.ElementCount == this.ArraySize {
		this.resize(this.ElementCount + 1)
	} else {
		this.ElementCount++
	}
	for i := this.ElementCount - 1; i > index; i-- {
		this.Array[i] = this.Array[i-1]
	}
}

func (this *Vector) increment() {
	if this.ElementCount == this.ArraySize {
		this.resize(this.ElementCount + 1)
	} else {
		this.ElementCount++
	}
}

func (this *Vector) decrement() {
	assert(this.ElementCount != 0, "Vector<T>::decrement - cannot decrement zero-length vector.")
	this.ElementCount--
}

func (this *Vector) resize(newCount int) {
	if newCount > 0 {
		blocks := newCount / VectorBlockSize
		if newCount%VectorBlockSize != 0 {
			blocks++
		}
		this.ElementCount = newCount
		this.ArraySize = blocks * VectorBlockSize
		newArray := make([]interface{}, this.ArraySize+1)
		copy(newArray, this.Array)
		this.Array = newArray
	}
}

func (this *Vector) Erase(index int) {
	assert(index < this.ElementCount, "Vector<T>::erase - out of bounds index.")
	if index < this.ElementCount-1 {
		copy(this.Array[index:this.ElementCount], this.Array[index+1:this.ElementCount])
	}
	this.ElementCount--
}

func (this *Vector) PushFront(value interface{}) {
	this.insert(0)
	this.Array[0] = value
}

func (this *Vector) PushBack(value interface{}) {
	this.increment()
	this.Array[this.ElementCount-1] = value
}

func (this *Vector) PopFront() {
	assert(this.ElementCount != 0, "Vector<T>::pop_front - cannot pop the front of a zero-length vector.")
	this.Erase(0)
}
func (this *Vector) PopBack() {
	assert(this.ElementCount != 0, "Vector<T>::pop_back - cannot pop the back of a zero-length vector.")
	this.decrement()
}

// Check that the index is within bounds of the list
func (this *Vector) withinRange(index int) bool {
	return index >= 0 && index < this.ElementCount
}

func (this *Vector) Back() interface{} {
	assert(this.ElementCount != 0, "Vector<T>::last - Error, no last element of a zero sized array! (const)")
	return this.Array[this.ElementCount-1]
}

func (this *Vector) Empty() bool {
	return this.ElementCount == 0
}

func (this *Vector) Size() int {
	return this.ArraySize
}

func (this *Vector) Clear() {
	this.ElementCount = 0
}

func (this *Vector) Len() int {
	return this.ElementCount
}

func (this *Vector) Get(index int) interface{} {
	assert(index < this.ElementCount, "Vector<T>::operator[] - out of bounds array access!")
	return this.Array[index]
}

func (this *Vector) Values() []interface{} {
	return this.Array[0:this.ElementCount]
}

func (this *Vector) Swap(i, j int) {
	this.Array[i], this.Array[j] = this.Array[j], this.Array[i]
}

func (this *Vector) Less(i, j int) bool {
	return true
}

func NewVector() *Vector {
	return &Vector{}
}
