package initialize

import (
	"github.com/bluele/gcache"
)

// InitGCache 初始化 GCache
func InitGCache() gcache.Cache {
	return gcache.New(256).ARC().Build()
}
