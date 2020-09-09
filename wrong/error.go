package wrong

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
)

// 定义错误
type Err struct {
	Code    int         `json:"status_code"` // 错误码
	Message string      `json:"message"`     // 展示给用户看的
	Errord  error       // 保存内部错误信息
	Debug   interface{} `json:"debug"`
	Trace   string      `json:"trace"`
	stack   *stack
}

func (err *Err) Error() string {
	return fmt.Sprintf("Err - code: %d, message: %s, error: %s", err.Code, err.Message, err.Errord)
}

// New returns an error that formats as the given text.
func New(code int, err error, message ...string) *Err {
	var msg string
	msg = ""
	if len(message) > 0 {
		msg = message[0]
	}
	return &Err{
		Code:    code,
		Message: msg,
		Errord:  err,
		stack:   callers(),
	}
}

// 解码错误, 获取 Code 和 Message
func DecodeErr(err error) (int, string) {
	if err == nil {
		return http.StatusOK, http.StatusText(http.StatusOK)
	}
	switch typed := err.(type) {
	case *Err:
		if typed.Code == http.StatusExpectationFailed {
			typed.Message = typed.Message + typed.Errord.Error()
		}
		return typed.Code, typed.Message
	default:
	}
	return http.StatusInternalServerError, err.Error()
}

func (err *Err) Format() string {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	var i int
	pc := *err.stack
	stackCount := len(*err.stack)
	for i = 0; i < stackCount; i++ { // Skip the expected number of frames
		f := runtime.FuncForPC(pc[i])
		file, line := f.FileLine(pc[i])

		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc[i])
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(f), source(lines, line))
	}
	return string(buf.Bytes()[:])
}

type stack []uintptr

func callers() *stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	var st stack = pcs[0:n]
	return &st
}

// function returns, if possible, the name of the function containing the PC.
func function(fn *runtime.Func) []byte {
	if fn == nil {
		return dunno()
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastSlash := bytes.LastIndex(name, slash()); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := bytes.Index(name, dot()); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot(), dot(), -1)
	return name
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno()
	}
	return bytes.TrimSpace(lines[n])
}
func slash() []byte {
	return []byte("/")
}
func dot() []byte {
	return []byte(".")
}
func centerDot() []byte {
	return []byte("·")
}
func dunno() []byte {
	return []byte("???")
}
