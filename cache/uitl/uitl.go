package uitl

import (
	"github.com/erDong01/micro-kit/cache/config"
	"strings"
	"time"
)

// 转换过期时间 -1 获取系统默认值
func ExpireDuration(ttl int) (exp time.Duration) {
	if -1 == ttl {
		return time.Duration(config.GetExpire()) * time.Second
	}
	return time.Duration(ttl) * time.Second
}

// 拼接key前缀
func PrefixKey(key string) string {
	return strings.Join([]string{config.GetPrefix(), ":", key}, "")
}
