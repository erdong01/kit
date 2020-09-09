package log

import (
	"github.com/sirupsen/logrus"
)

// 日志输出字段
type JSON logrus.Fields

// 日志输出组件
var logger *logrus.Logger

func init() {
	logger = NewLogrus(getConfig())
}

// 信息
func Info(args ...interface{}) {
	logger.Info(args...)
}

// 警告
func Warn(args ...interface{}) {
	logger.Warn(args...)
}

// 错误
func Error(args ...interface{}) {
	logger.Error(args...)
}

// 退出程序
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

// 抛出异常
func Panic(args ...interface{}) {
	logger.Panic(args...)
}

// 信息 输出键值对
func InfoFields(f map[string]interface{}) {
	logger.WithFields(f).Info("")
}

// 警告 输出键值对
func WarnFields(f map[string]interface{}) {
	logger.WithFields(f).Warn("")
}

// 错误 输出键值对
func ErrorFields(f map[string]interface{}) {
	logger.WithFields(f).Error("")
}

// 退出程序 输出键值对
func FatalFields(f map[string]interface{}) {
	logger.WithFields(f).Fatal("")
}

// 异常 输出键值对
func PanicFields(f map[string]interface{}) {
	logger.WithFields(f).Panic("")
}
// 打印
func Print(args ...interface{}) {
	logger.Print(args...)
}

