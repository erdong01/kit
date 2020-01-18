package cache

import (
	"rxt/internal/cache/I"
	"rxt/internal/core"
)

func New() *I.ICache {
	return core.New().Cache
}
