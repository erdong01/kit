package log

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/sirupsen/logrus"
)

// path 日志目录
// logName 日志名称
func NewLogrus(cacheConfig *Config) *logrus.Logger {

	l := logrus.New()

	outWriter := strings.Split(cacheConfig.Writer, ",")

	// 设置输出格式
	if "json" == cacheConfig.Format {
		l.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	// 设置输出等级
	level, err := logrus.ParseLevel(cacheConfig.Level)
	if err != nil {
		l.Fatal("日志登记错误:", err)
	}
	l.SetLevel(level)

	// 设置日志输出位置
	switch len(outWriter) {
	case 1:
		if "" == outWriter[0] {
			l.SetOutput(ioutil.Discard)
		} else if "stderr" == outWriter[0] {
			l.SetOutput(os.Stderr)
		} else if "file" == outWriter[0] {
			fileWriter, err := writerFile(cacheConfig)
			if err != nil {
				l.Fatal("创建日志文件失败:", err)
			}
			l.SetOutput(fileWriter)
		}
	case 2:
		fileWriter, err := writerFile(cacheConfig)
		if err != nil {
			l.Fatal("创建日志文件失败:", err)
		}
		l.SetOutput(io.MultiWriter(os.Stderr, fileWriter))
	default:
		l.Fatal("日志输出位置错误")
	}
	return l
}

// 判断/创建文件夹
func isExistOrCreate(path string) {
	_, err := os.Stat(path)
	if err != nil && !os.IsExist(err) {
		os.Mkdir(path, 0755)
	}
}

func writerFile(cacheConfig *Config) (io.WriteCloser, error) {
	isExistOrCreate(cacheConfig.File.Path)
	logFile := path.Join(cacheConfig.File.Path, cacheConfig.File.Name)
	// 是否开启日志分割
	if cacheConfig.File.RotationOpen {
		var logUnit time.Duration
		var logTimeFormat string
		switch cacheConfig.File.RotationUnit {
		case "minute":
			logUnit = 1
			logTimeFormat = "%Y-%m-%d-%H-%M"
		case "hour":
			logUnit = 60
			logTimeFormat = "%Y-%m-%d-%H"
		case "day":
			logUnit = 1440
			logTimeFormat = "%Y-%m-%d"
		default:
			return nil, errors.New("日志分割时间单位错误")
		}
		writer, err := rotatelogs.New(
			logFile+"."+logTimeFormat+".log",
			// 生成软链，指向最新日志文件
			//rotatelogs.WithLinkName(logFile),
			// 文件最大保存时间
			rotatelogs.WithMaxAge(time.Minute*logUnit*time.Duration(cacheConfig.File.RotationTime)),
			// 日志切割时间间隔
			rotatelogs.WithRotationTime(time.Minute*logUnit*time.Duration(cacheConfig.File.rotationTimeSave)),
		)
		return writer, err
	} else {
		return os.OpenFile(logFile+".log", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	}

}

// 日志切割
// logPath 日志目录 logFileName 日志文件名称 maxAge 文件最大保存时间 rotationTime 日志切割时间间隔
//func NewIfsHook(logPath string, logFileName string, maxAge time.Duration, rotationTime time.Duration) logrus.Hook {
//	baseLogPaht := path.Join(logPath, logFileName)
//	writer, err := rotatelogs.New(
//		baseLogPaht+".%Y%m%d%H%M.log",
//		//baseLogPaht+"-%Y-%m-%d.log",
//		rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
//		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
//	)
//	if err != nil {
//		Error("config local file system logger error:" + errors.WithStack(err).Error())
//	}
//	return lfshook.NewHook(lfshook.WriterMap{
//		logrus.DebugLevel: writer, // 为不同级别设置不同的输出目的
//		logrus.InfoLevel:  writer,
//		logrus.WarnLevel:  writer,
//		logrus.ErrorLevel: writer,
//		logrus.FatalLevel: writer,
//		logrus.PanicLevel: writer,
//	}, &logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"})
//}

//// 绑定文件输出
//func bindFile() (err error) {
//	logFile := filepath.Join(f.fileDir, f.fileName)
//	f.logFile, err = os.OpenFile(logFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
//	if err != nil {
//		return
//	}
//	f.lg.SetOutput(f.logFile)
//	return
//}
