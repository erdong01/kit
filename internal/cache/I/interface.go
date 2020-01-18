package I

// 接口实现必须是并发安全
type Cache interface {

	// 判断是否有缓存
	// key 缓存名称
	Has(key string) bool

	// 设置基础数据缓存
	// key 缓存名称
	// value 缓存内容
	// expire 缓存自动过期时间
	Set(key string, value interface{}, expire int) error

	// 获取基础数据缓存
	// key 缓存名称
	// value 指针类型
	Get(key string, value interface{}) error

	// 设置键值对数据缓存
	// key 缓存名称
	// expire 缓存自动过期时间
	SetJSON(key string, value interface{}, expire int) error

	// 获取键值对数据缓存
	// key 缓存名称
	// value 指针类型
	GetJSON(key string, value interface{}) error

	// 获取数据
	// key 缓存名称
	// 返回byte切片
	GetBytes(key string) ([]byte, error)

	// 删除缓存
	// key 缓存名称
	Del(key string) error
}

type ICache struct {
	Cache
}
