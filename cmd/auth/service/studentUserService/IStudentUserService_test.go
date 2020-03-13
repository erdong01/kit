package studentUserService

import (
	"fmt"
	"rxt/internal/core/register"
	"testing"
)

func TestV1_Login(t *testing.T) {
	register.GlobalInit()
	var param = Param{
		StudentUserLoginName: "70136610",
		Password:             "70136610",
	}
	res, err := New().Login(&param)
	fmt.Println(res, err)
}
