package config

// 获取缓存驱动
func GetDrive() string {
	return "redis"
}

// 获取自动过期时间
func GetExpire() int {
	return 30
}

// 获取前缀
func GetPrefix() string {
	return "api_517rxt_com"
}
