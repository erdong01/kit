package log

import (
	"os"
	"os/exec"
	"testing"
)

func TestMain(m *testing.M) {
	m.Run()
	// TODO 删除log文件
}

func TestInfo(t *testing.T) {
	Info("info")
}

func TestWarn(t *testing.T) {
	Warn("warn")
}

func TestError(t *testing.T) {
	Error("Error")
}

func TestFatal(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		Fatal("Fatal")
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestFatal")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Log("程序未退出")
	t.Fail()
}

func TestPanic(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Fail()
		}
	}()
	Panic("panic")
}

func TestInfoFields(t *testing.T) {
	InfoFields(JSON{
		"message": "info",
	})
}

func TestFatalFields(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		Fatal("Fatal")
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestFatalFields")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Log("程序未退出")
	t.Fail()
}

func TestPanicFields(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Fail()
		}
	}()
	PanicFields(JSON{
		"message": "Panic",
	})
}

func BenchmarkInfo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Info(i)
	}
}
