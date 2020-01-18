package cache

import (
	"rxt/internal/cache/I"
	"rxt/internal/cache/config"
	"rxt/internal/cache/uitl"
	"rxt/internal/db/redis"
)

const (
	DefaultTTl int = -1 // 系统默认自动过期时间值
	ForeverTTl int = 0  // 永不过期
)

type JSON map[string]interface{}

func make() I.Cache {
	drive := config.GetDrive()
	if drive == "redis" {
		return redis.NewCache()
	} else {
		panic("驱动不存在")
	}
}

// 判断缓存是否存在
// true 存在
// false 不存在
func Has(key string) bool {
	return make().Has(key)
}

// 获取 key 值对应基础类型数据
// value 用于赋值结果的指针类型
func Get(key string, value interface{}) error {
	return make().Get(key, value)
}

// 获取 key 值对应的键值对数据 如: map []map struct []struct
// value 用于赋值结果的指针类型
func GetJSON(key string, value interface{}) error {
	return make().GetJSON(key, value)
}

// 设置 key 对应的数据
// value 值内容为基础类型数据 如: string int []byte 等
// expire 自动过期时间 单位秒
func Set(key string, value interface{}, expire int) error {
	return make().Set(key, value, expire)
}

// 设置 key 对应的数据
// value 值内容为键值对 如: map []map struct []struct 等
// expire 自动过期时间 单位秒
func SetJSON(key string, value interface{}, expire int) error {
	return make().SetJSON(uitl.PrefixKey(key), value, expire)
}

// 返回key值对应的byte切片数据
func GetBytes(key string) ([]byte, error) {
	return make().GetBytes(uitl.PrefixKey(key))
}

// 删除 key 对应的值
func Del(key string) error {
	return make().Del(uitl.PrefixKey(key))
}
