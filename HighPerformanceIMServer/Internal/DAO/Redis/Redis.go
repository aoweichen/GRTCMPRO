package Redis

import (
	"HighPerformanceIMServer/Configs"
	"context"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

// DataBase REDISDB 为全局变量，表示Redis客户端
var DataBase *redis.Client

// InitClient 函数用于初始化Redis客户端连接
func InitClient() {
	// 创建一个新的Redis客户端实例，并配置相关参数
	DataBase = redis.NewClient(&redis.Options{
		Network:      "tcp",                                                               // 使用TCP网络连接
		Addr:         Configs.ConfigData.Redis.Host + ":" + Configs.ConfigData.Redis.Port, // Redis服务器的地址和端口号
		Password:     Configs.ConfigData.Redis.Password,                                   // Redis服务器的密码
		DB:           Configs.ConfigData.Redis.DB,                                         // 要连接的Redis数据库编号
		PoolSize:     Configs.ConfigData.Redis.Poll,                                       // 连接池大小，默认为CPU核心数的四倍
		MinIdleConns: Configs.ConfigData.Redis.Conn,                                       // 在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量
		DialTimeout:  5 * time.Second,                                                     // 连接超时时间
		ReadTimeout:  5 * time.Second,                                                     // 读取超时时间
		WriteTimeout: 5 * time.Second,                                                     // 写入超时时间
		PoolTimeout:  5 * time.Second,                                                     // 连接池超时时间
	})

	// 发送Ping命令到Redis服务器，检查连接是否正常
	_, err := DataBase.Ping(context.Background()).Result()
	// 如果连接出错，则记录错误日志并退出程序
	if err != nil {
		zap.S().Errorln(err)
	}
}
