package cache

import (
	"github.com/erDong01/micro-kit/core"
	"github.com/erDong01/micro-kit/internal/cache/I"
)

func New() *I.ICache {
	return core.New().Cache
}
