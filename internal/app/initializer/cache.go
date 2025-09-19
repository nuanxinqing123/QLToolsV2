package initializer

import (
	"github.com/bluele/gcache"
)

func Cache() gcache.Cache {
	return gcache.New(256).ARC().Build()
}
