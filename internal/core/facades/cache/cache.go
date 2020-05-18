package cache

import (
	"github.com/erDong01/micro-kit/internal/cache/I"
	"github.com/erDong01/micro-kit/internal/core"
)

func New() *I.ICache {
	return core.New().Cache
}
