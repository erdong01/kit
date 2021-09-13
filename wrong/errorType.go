package wrong

import (
	"fmt"
	"log"
	"runtime"
)

const (
	NONE_ERROR      = iota
	VERSION_ERROR   //版本不正确
	ACCOUNT_NOEXIST //账号不存在
	PASSWORD_ERROR  //密码不正确
)

//TraceCode 输出错误，跟踪代码
func TraceCode(code ...interface{}) {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	data := ""
	for _, v := range code {
		data += fmt.Sprintf("%v", v)
	}
	data += string(buf[:n])
	log.Printf("==> %s\n", data)
}
