package configs

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type ServerConfig struct {
	LoadBalanceService LoadBalanceServiceConfig
	Redis              RedisConfig
	Logger             LogConfig
	Consul             ConsulConfigs
	AuthService        AuthServiceConfigs
	IMServers          IMServersConfigs
}

type LoadBalanceServiceConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type RedisConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

// LogConfig LogConf 结构体用于存储日志相关的配置信息
type LogConfig struct {
	Level     string `json:"level"`     // 日志级别
	Type      string `json:"type"`      // 日志类型
	FileName  string `json:"filename"`  // 日志文件名
	MaxSize   int    `json:"maxSize"`   // 日志文件最大大小（字节）
	MaxBackup int    `json:"maxBackup"` // 日志文件最大备份数
	MaxAge    int    `json:"maxAge"`    // 日志文件最大保存时间（天）
	Compress  bool   `json:"compress"`  // 是否启用日志压缩
}

type ConsulConfigs struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type AuthServiceConfigs struct {
	Name string `json:"name"`
}

type IMServersConfigs struct {
	Name string `json:"name"`
}

var ConfigData ServerConfig

func InitConfigs() {
	// 配置文件名和路径
	viper.SetConfigFile("configs.yaml")
	// 设置配置文件类型
	viper.SetConfigType("yaml")
	// 监视配置文件变化
	viper.WatchConfig()

	// 读取稳健配置项
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		zap.S().Error("viper ReadInConfig Error: ", err.Error())
		panic(err.Error())
	}
	// 写变量
	err := viper.Unmarshal(&ConfigData)
	if err != nil {
		fmt.Println(err)
		zap.S().Infoln("viper Unmarshal ConfigData Error:", err.Error())
		panic(err.Error())
	}

	// 监听配置文件的变化并动态更新
	viper.OnConfigChange(func(e fsnotify.Event) {
		zap.S().Infoln("Config file changed:", e.Name)
	})
}
