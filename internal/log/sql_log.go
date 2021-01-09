/**
记录sql日志
*/
package log

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"io"
	"log"
	"os"
	"time"
)

var Sl *log.Logger

type sqlLog struct {
}

func init() {
	_, err := os.Stat("log")
	if os.IsNotExist(err) {
		os.Mkdir("log", 0666)
	}

	fileName := fmt.Sprintf("log/sql-%v.log", time.Now().Format("2006-01-02"))
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	Sl = log.New(io.MultiWriter(os.Stderr, file), "", log.LUTC)
}

func (logger *sqlLog) Print(values ...interface{}) {
	sqlStr := gorm.LogFormatter(values...)[3]
	currentTime := fmt.Sprintf("[%v] ", time.Now().Format("2006-01-02 15:04:05"))
	filePath := fmt.Sprintf(" (%v)", values[1])
	executeTime := fmt.Sprintf("[%v] ", values[2])
	logStr := currentTime + executeTime + sqlStr.(string) + filePath
	Sl.Println(logStr)
}

// 开启 gorm 日志，并记录到logs文件夹下。
func SetSqlLogger(db *gorm.DB) {
	// TODO 生产环境判断和慢查询时间设定
	db.LogMode(true)
	sqlLog := &sqlLog{}
	db.SetLogger(sqlLog)
}
