package drive

import (
	_ "github.com/go-sql-driver/mysql" // 引入mysql驱动

	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/erDong01/gin-kit/internal/config"
	"github.com/erDong01/gin-kit/internal/log"
)

// New 初始化数据库
func New() *gorm.DB {
	mysqlCnf := config.GetMySQL()
	if err := config.New().Get(&mysqlCnf, "mysql"); err != nil {
		log.Fatal(err)
	}

	db, err := gorm.Open("mysql", DSN(mysqlCnf))
	fmt.Println(mysqlCnf)
	if err != nil {
		log.Fatal(err)
	}

	db.DB().SetMaxIdleConns(mysqlCnf.MaxIdleConn)
	db.DB().SetMaxOpenConns(mysqlCnf.MaxOpenConn)
	db.LogMode(true)
	return db
}

// DSN 数据库连接串
func DSN(c *config.MySQL) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s&charset=utf8&parseTime=true",
		c.User, c.Password, c.Host, c.Port, c.Database, c.Parameters)
}
