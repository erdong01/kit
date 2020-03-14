package aiLogic

import (
	"rxt/cmd/exam/logic/base"
)

type V1 struct {
	base.Logic
}

func (c V1) Dina(param string) (string, error) {
	return "调用成功", nil
}
