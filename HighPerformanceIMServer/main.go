package main

import (
	"HighPerformanceIMServer/Configs"
	"HighPerformanceIMServer/Internal/Api/Services/Clients/Manager"
	"HighPerformanceIMServer/Internal/Api/Services/GRPCService"
	"HighPerformanceIMServer/Internal/DAO/MYSQL"
	"HighPerformanceIMServer/Internal/Middleware"
	"HighPerformanceIMServer/Internal/Router"
	"HighPerformanceIMServer/Packages/Consul"
	"HighPerformanceIMServer/Packages/CoroutinesPoll"
	"HighPerformanceIMServer/Packages/Logger"
	"HighPerformanceIMServer/Packages/MessageQueue/NSQQueue"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitConfigs() {
	Configs.InitConfigs()                    // 初始化配置，从指定的配置文件中加载配置
	Logger.InitLogger()                      // 初始化日志记录器
	MYSQL.InitMySQLDB()                      // 初始化 MySQL 数据库连接 	// 初始化 Redis 客户端
	CoroutinesPoll.ConnectPool()             // 连接协程池
	err := NSQQueue.InitNewNSQProducerPoll() // 初始化 NSQ 生产者连接池
	if err != nil {
		panic(err) // 如果初始化失败，则抛出异常
	}
}
func StartService(router *gin.Engine) {
	InitConfigs()
	go Manager.IMMessageClientManager.StartServer() // 在后台启动 IMMessageClientManager 的服务器
	router.Use(Middleware.Recover)                  // 使用 Recover 中间件，用于处理请求发生的 panic
	SetRoute(router)                                // 设置路
	GRPCService.StartGRPCServer()                   // 启动 gRPC 服务器
	//
	service := &Consul.Service{
		ID:   Configs.ConfigData.Server.Name + Configs.ConfigData.Server.Node,
		Name: Configs.ConfigData.Server.Name,
		Host: Configs.ConfigData.Server.Host,
		Port: Configs.ConfigData.Server.ListenPort,
	}
	err := Consul.RegisterService(service, Configs.ConfigData.Consul.Host+Configs.ConfigData.Consul.Port)
	if err != nil {
		zap.S().Errorln(err)
		return
	}
	err = router.Run(":" + strconv.Itoa(Configs.ConfigData.Server.ListenPort)) // 运行 Gin 引擎，监听指定地址和端口
	if err != nil {
		panic(err) // 如果运行出错，则抛出异常
	}
}

func SetRoute(router *gin.Engine) {
	// 健康检查路由
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})
	Router.RegisterApiRoutes(router) // 注册 API 路由
	Router.RegisterWSRouters(router) // 注册 WebSocket 路由
}

func main() {

	r := gin.Default() // 创建一个默认的 Gin 引擎实例
	StartService(r)    // 启动服务
}
