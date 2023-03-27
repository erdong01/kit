package slice

import (
	"reflect"
)

func Del[T comparable](s *[]T, val T) {
	temp := *s
	if len(temp) == 0 {
		return
	}
	var i int
	for key, v := range temp {
		if v != val {
			temp[i] = temp[key]
			i++
		}
	}
	*s = temp[:i]
}

func DelFunc[T comparable](s *[]T, f func(i int) bool) {
	temp := *s
	if len(temp) == 0 {
		return
	}
	var i int
	for key := range temp {
		if !f(key) {
			temp[i] = temp[key]
			i++
		}
	}
	*s = temp[:i]
}

func DelByIndex[T any](s *[]T, index int) {
	temp := *s
	count := len(temp)
	if count == 0 || index > count {
		return
	}
	count--
	var i int = index
	for index < count {
		index++
		temp[i] = temp[index]
		i++
	}
	*s = temp[:i]
}

func Unique[T comparable](s *[]T) {
	var temp = *s
	if len(*s) == 0 {
		return
	}
	va := reflect.ValueOf(temp)
	var k int = 1
	for i := 1; i < va.Len(); i++ {
		if !reflect.DeepEqual(va.Index(i-1).Interface(), va.Index(i).Interface()) {
			temp[k] = temp[i]
			k++
		}
	}
	*s = temp[:k]
}

func Prepend[T any](s *[]T, val T) {
	*s = append([]T{val}, *s...)
}

func Insert[T any](s *[]T, index int, val T) {
	var temp = *s
	temp = append(temp, val)
	copy(temp[index+1:], temp[index:])
	temp[index] = val
	*s = temp
}

func Diff[T comparable](a, b []T) (deff []any) {
	bm := make(map[T]struct{}, len(b))
	for k := range b {
		bm[b[k]] = struct{}{}
	}
	for k := range a {
		if _, ok := bm[a[k]]; !ok {
			deff = append(deff, a[k])
		}
	}
	return
}

func Equal[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	if (a == nil) != (b == nil) {
		return false
	}
	b = b[:len(a)]
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
