package Configs

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type IMServerConfigs struct {
	Server ServerConfigs
	MySQL  MySQLConfigs
	Logger LogConfigs
	Nsq    NSQConfigs
	Consul ConsulConfigs
	QNY    QNYConfigs
	Redis  RedisConfig
}

type ServerConfigs struct {
	Name          string `json:"name"`
	Host          string `json:"host"`
	ListenPort    int    `json:"listen_port"`
	CoroutinePoll int    `json:"coroutine_poll"`
	FilePath      string `json:"file_path"`
	ClusterOpen   bool   `json:"cluster_open"`
	Node          string `json:"node"`
	GrpcListen    string `json:"grpcListen"` // gRPC 服务的监听地址，指定 gRPC 服务监听的 IP 和端口号
}

type MySQLConfigs struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
	Charset  string `json:"charset"`
}

type LogConfigs struct {
	Level       string `json:"level"`
	Type        string `json:"type"`
	LogFilePath string `json:"log_file_path"`
	MaxSize     int64  `json:"max_size"`
	MaxBackup   int64  `json:"max_backup"`
	MaxAge      int64  `json:"max_age"`
	Compress    bool   `json:"compress"`
}

// RedisConfig RedisConf 结构体用于存储 Redis 数据库的配置信息
type RedisConfig struct {
	Host     string `json:"host"`     // Redis 服务器主机地址
	Port     string `json:"port"`     // Redis 服务器端口号
	Password string `json:"password"` // Redis 服务器密码
	DB       int    `json:"db"`       // 要连接的 Redis 数据库索引
	Poll     int    `json:"poll"`     // Redis 服务器轮询间隔时间（毫秒）
	Conn     int    `json:"conn"`     // Redis 服务器最大连接数
}
type NSQConfigs struct {
	LookupHost string `json:"lookup_host"`
	NsqHost    string `json:"nsq_host"`
}

type ConsulConfigs struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type QNYConfigs struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Bucket    string `json:"bucket"`
	Domain    string `json:"domain"`
}

var ConfigData IMServerConfigs

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
