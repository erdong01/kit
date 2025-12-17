package sharding

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/erdong01/sharding"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DBPool []*gorm.DB

// ---------------------------------------------------------
// 1. 基因法相关 (保持不变)
// ---------------------------------------------------------
func GenerateOrderNo(userID int64) string {
	timestamp := time.Now().UnixMilli()
	gene := userID % 10000
	random := rand.Intn(100)
	return fmt.Sprintf("%d%04d%02d", timestamp, gene, random)
}

// ---------------------------------------------------------
// 2. 自定义分表算法 (Key -> 0~1023)
// ---------------------------------------------------------
func MeituanTableSharding(value interface{}) (suffix string, err error) {
	var routingKey int64
	// ... 解析 UserID 或 OrderNo 的逻辑与之前一致 ...
	switch v := value.(type) {
	case int64:
		routingKey = v
	case string: // 解析 OrderNo 里的基因
		if len(v) < 6 {
			return "", errors.New("invalid order_no")
		}
		geneStr := v[len(v)-6 : len(v)-2]
		geneInt, _ := strconv.ParseInt(geneStr, 10, 64)
		routingKey = geneInt
	default:
		return "", errors.New("unsupported key")
	}

	// 算出 0-1023 的表后缀
	// 注意：虽然物理上分了库，但逻辑上我们还是认为有 1024 张唯一的表名
	// 比如 orders_0000 在 db_0, orders_0032 在 db_1
	shardingIndex := routingKey % 1024
	return fmt.Sprintf("_%04d", shardingIndex), nil
}

// 核心 分库路由选择器
func GetDB(value any) *gorm.DB {
	var routingKey int64
	switch v := value.(type) {
	case int64:
		routingKey = v
	case string:
		if len(v) >= 6 {
			geneStr := v[len(v)-6 : len(v)-2]
			routingKey, _ = strconv.ParseInt(geneStr, 10, 64)
		}
	}
	shardingIndex := routingKey % 1024
	dbIndex := shardingIndex / 32
	if dbIndex >= 32 {
		dbIndex = 0 // 防御性代码
	}
	return DBPool[dbIndex]
}

func Init() {
	baseDSN := "root:123456@tcp(127.0.0.1:3306)/my_sharding_db_%02d?charset=utf8mb4&parseTime=True&loc=Local"

	for i := 0; i < 32; i++ {
		dsn := fmt.Sprintf(baseDSN, i)
		// 实际项目中，这里建议复用 mysql.Open 的配置，避免过多连接
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(fmt.Sprintf("连接库 db_%02d 失败: %v", i, err))
		}
		middleware := sharding.Register(sharding.Config{
			ShardingKey:         "user_id",
			NumberOfShards:      1024, // 告诉插件总共有 1024 张表逻辑
			PrimaryKeyGenerator: sharding.PKSnowflake,
			ShardingAlgorithm:   MeituanTableSharding,
		}) // *** 修正点：这里传入要分表的 Model ***

		db.Use(middleware)

		DBPool[i] = db
		fmt.Printf("DB_%02d 初始化完成\n", i)
	}
}
