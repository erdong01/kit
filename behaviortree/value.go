package behaviortree

import (
	"errors"
	"reflect"
	"runtime"
	"sync"
)

var (
	runtimeCallers       = runtime.Callers
	runtimeCallersFrames = runtime.CallersFrames
	runtimeFuncForPC     = runtime.FuncForPC
)

var (
	valueCallMutex  sync.Mutex
	valueDataMutex  sync.RWMutex
	valueDataKey    interface{}
	valueDataChan   chan interface{}
	valueDataCaller [1]uintptr
)

func (n Node) WithValue(key, value interface{}) Node {
	if n == nil {
		panic(errors.New(`behaviortree.Node.WithValue nil receiver`))
	}
	if key == nil {
		panic(errors.New(`behaviortree.Node.WithValue nil key`))
	}
	if !reflect.TypeOf(key).Comparable() {
		panic(errors.New(`behaviortree.Node.WithValue key is not comparable`))
	}
	return func() (Tick, []Node) {
		n.valueHandle(func(k interface{}) (interface{}, bool) {
			if k == key {
				return value, true
			}
			return nil, false
		})
		return n()
	}
}

func (n Node) Value(key interface{}) interface{} {
	if n != nil {
		valueCallMutex.Lock()
		defer valueCallMutex.Unlock()
		return n.valueSync(key)
	}
	return nil
}

func (n Node) valueSync(key interface{}) (value interface{}) {
	if n.valuePrep(key) {
		select {
		case value = <-valueDataChan:
		default:
		}
		valueDataMutex.Lock()
		valueDataKey = nil
		valueDataChan = nil
		valueDataMutex.Unlock()
	}
	return
}

func (n Node) valuePrep(key interface{}) bool {
	valueDataMutex.Lock()
	if runtimeCallers(2, valueDataCaller[:]) < 1 {
		valueDataMutex.Unlock()
		return false
	}
	valueDataKey = key
	valueDataChan = make(chan interface{}, 1)
	valueDataMutex.Unlock()
	n()
	return true
}

func (n Node) valueHandle(fn func(key interface{}) (interface{}, bool)) {
	valueDataMutex.RLock()
	dataKey, dataChan, dataCaller := valueDataKey, valueDataChan, valueDataCaller
	valueDataMutex.RUnlock()
	value, ok := fn(dataKey)
	if !ok {
		return
	}
	dataKey = nil
	const depth = 2 << 7
	callers := make([]uintptr, depth)
	for skip := 4; skip > 0; skip += depth {
		callers = callers[:runtimeCallers(skip, callers[:])]
		for _, caller := range callers {
			if caller == dataCaller[0] {
				select {
				case dataChan <- value:
				default:
				}
				return
			}
		}
		if len(callers) != depth {
			return
		}
	}
}
