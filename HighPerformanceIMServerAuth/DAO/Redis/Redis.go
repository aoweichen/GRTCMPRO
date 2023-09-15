package Redis

import (
	"HighPerformanceIMServerAuth/Configs"
	"context"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

var AuthRedisDB *redis.Client

func InitRedisDB() {
	AuthRedisDB = redis.NewClient(&redis.Options{
		Network:      "tcp",
		Addr:         Configs.ConfigData.Redis.Host + ":" + Configs.ConfigData.Redis.Port,
		Password:     Configs.ConfigData.Redis.Password,
		DB:           int(Configs.ConfigData.Redis.Database),
		PoolSize:     int(Configs.ConfigData.Redis.Poll), // 连接池大小，默认为CPU核心数的四倍
		MinIdleConns: int(Configs.ConfigData.Redis.Conn), // 在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量
		DialTimeout:  5 * time.Second,                    // 连接超时时间
		ReadTimeout:  5 * time.Second,                    // 读取超时时间
		WriteTimeout: 5 * time.Second,                    // 写入超时时间
		PoolTimeout:  5 * time.Second,
	})
	// 发送Ping命令到Redis服务器，检查连接是否正常
	_, err := AuthRedisDB.Ping(context.Background()).Result()
	// 如果连接出错，则记录错误日志并退出程序
	if err != nil {
		zap.S().Errorln(err)
	}
}
