package Configs

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type ServerConfig struct {
	AuthService AuthGRPCServiceConfig
	MySQL       MySQLConfig
	Redis       RedisConfig
	Logger      LogConfig
	JWT         JWTConfig
	Mail        MailConfig
	CONSUL      ConsulConfig
}

type AuthGRPCServiceConfig struct {
	Name       string `json:"name"`
	Host       string `json:"host"`
	ListenPort int    `json:"listen_port"`
}

type MySQLConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
	Charset  string `json:"charset"`
}

type LogConfig struct {
	Level       string `json:"level"`
	Type        string `json:"type"`
	LogFilePath string `json:"log_file_path"`
	MaxSize     int64  `json:"max_size"`
	MaxBackup   int64  `json:"max_backup"`
	MaxAge      int64  `json:"max_age"`
	Compress    bool   `json:"compress"`
}

type JWTConfig struct {
	Secret          string `json:"secret"`
	TokenTimeToLive int64  `json:"token_time_to_live"`
}

type RedisConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
	Database int64  `json:"database"`
	Poll     int64  `json:"poll"`
	Conn     int64  `json:"conn"`
}

type MailConfig struct {
	Driver                        string `json:"driver"`
	Host                          string `json:"host"`
	Name                          string `json:"name"`
	Password                      string `json:"password"`
	Port                          int64  `json:"port"`
	Encryption                    string `json:"encryption"`
	FromName                      string `json:"from_name"`
	EmailCodeSubject              string `json:"email_code_subject"`
	EmailCodeHtmlTemplateFilePath string `json:"email_code_html_template_file_path"`
}

type ConsulConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
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
