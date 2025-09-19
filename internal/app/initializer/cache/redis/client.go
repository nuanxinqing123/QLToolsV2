package redis

import (
	"context"
	"fmt"

	"github.com/nuanxinqing123/QLToolsV2/internal/app/config"
	"github.com/redis/go-redis/v9"
)

// CacheRedis 初始化Redis缓存
func CacheRedis() *redis.Client {
	_cache := config.Config.Cache
	RDB := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", _cache.Host, _cache.Port),
		Password: _cache.Password, // no password set
		DB:       _cache.DB,       // use default DB
		PoolSize: _cache.PoolSize, // 连接池大小
	})

	// 测试连接
	_, err := RDB.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}

	return RDB
}
