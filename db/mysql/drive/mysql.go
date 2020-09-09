package drive

import (
	"fmt"
	"github.com/erDong01/micro-kit/internal/config"
	"github.com/erDong01/micro-kit/internal/log"
	_ "github.com/go-sql-driver/mysql" // 引入mysql驱动
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

// New 初始化数据库
func New() *gorm.DB {
	mysqlCnf := config.GetMySQL()
	if err := config.New().Get(&mysqlCnf, "mysql"); err != nil {
		log.Fatal(err)
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		DriverName: "my_mysql_driver",
		DSN:        DSN(mysqlCnf),
	}), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	mySqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	mySqlDB.SetMaxIdleConns(mysqlCnf.MaxIdleConn)
	mySqlDB.SetMaxOpenConns(mysqlCnf.MaxOpenConn)
	mySqlDB.SetConnMaxLifetime(time.Hour)
	return db
}

// DSN 数据库连接串
func DSN(c *config.MySQL) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s&charset=utf8&parseTime=true",
		c.User, c.Password, c.Host, c.Port, c.Database, c.Parameters)
}
