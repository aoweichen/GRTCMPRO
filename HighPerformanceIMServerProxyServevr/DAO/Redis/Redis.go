package Redis

import (
	"HighPerformanceIMServerProxyServevr/configs"
	"github.com/redis/go-redis/v9"
)

var ProxyRedis *redis.Client

func InitRedisDB() {
	ProxyRedis = redis.NewClient(&redis.Options{
		Addr:     configs.ConfigData.Redis.Host + configs.ConfigData.Redis.Port,
		Password: configs.ConfigData.Redis.Password,
		DB:       configs.ConfigData.Redis.DB,
	})
}
