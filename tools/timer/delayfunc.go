package timer

import (
	"fmt"
	"log"
	"reflect"
)

/*
   定义一个延迟调用函数
	延迟调用函数就是 时间定时器超时的时候，触发的事先注册好的
	回调函数
*/

// DelayFunc 延迟调用函数对象
type DelayFunc struct {
	f    func(...interface{}) //f : 延迟函数调用原型
	args []interface{}        //args: 延迟调用函数传递的形参
}

// NewDelayFunc 创建一个延迟调用函数
func NewDelayFunc(f func(v ...interface{}), args []interface{}) *DelayFunc {
	return &DelayFunc{
		f:    f,
		args: args,
	}
}

// String 打印当前延迟函数的信息，用于日志记录
func (df *DelayFunc) String() string {
	return fmt.Sprintf("{DelayFun:%s, args:%v}", reflect.TypeOf(df.f).Name(), df.args)
}

func (df *DelayFunc) Call() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(df.String(), "Call err:", err)
		}
	}()

	df.f(df.args...)
}
