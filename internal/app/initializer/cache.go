package initializer

import (
	_redis "github.com/nuanxinqing123/QLToolsV2/internal/app/initializer/cache/redis"
	"github.com/redis/go-redis/v9"
)

func Cache() *redis.Client {
	return _redis.CacheRedis()
}
