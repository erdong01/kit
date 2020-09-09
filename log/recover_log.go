package log

import "io"

// 系统错误处理输出至日志
type recoverLog struct {
}

func (r *recoverLog) Write(p []byte) (n int, err error) {
	Error(string(p))
	return len(p), nil
}

// 返回输入结构体
func GetRecoverLog() io.Writer {
	return &recoverLog{}
}
