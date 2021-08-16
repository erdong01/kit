package behaviortree

import "reflect"

type (
	Frame struct {
		PC uintptr

		Function string

		File string

		Line int

		Entry uintptr
	}
	vkFrame struct{}
)

func (n Node) Frame() *Frame {
	if v, _ := n.Value(vkFrame{}).(*Frame); v != nil {
		v := *v
		return &v
	}
	return newFrame(n)
}

func (t Tick) Frame() *Frame { return newFrame(t) }

func newFrame(v interface{}) (f *Frame) {
	if v := reflect.ValueOf(v); v.IsValid() && v.Kind() == reflect.Func && !v.IsNil() {
		p := v.Pointer()
		if v := runtimeFuncForPC(p); v != nil {
			f = &Frame{
				PC:       p,
				Function: v.Name(),
				Entry:    v.Entry(),
			}
			f.File, f.Line = v.FileLine(f.Entry)
		}
	}
	return
}
