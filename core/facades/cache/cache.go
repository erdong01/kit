package cache

import (
	"github.com/erDong01/micro-kit/cache/I"
	"github.com/erDong01/micro-kit/core"
)

func New() *I.ICache {
	return core.New().Cache
}
