package main

import (
	"HighPerformanceIMServerProxyServevr/DAO/Redis"
	"HighPerformanceIMServerProxyServevr/GlobalVars"
	"HighPerformanceIMServerProxyServevr/Logger"
	"HighPerformanceIMServerProxyServevr/Route"
	"HighPerformanceIMServerProxyServevr/configs"
	"HighPerformanceIMServerProxyServevr/packages/LoadBalanceHTTP"

	"github.com/gin-gonic/gin"
)

func StartLoadBalanceService() {
	configs.InitConfigs()
	GlobalVars.InitGlobalVars()
	Logger.InitLogger()
	Redis.InitRedisDB()
	LoadBalanceHTTP.InitializeResponseTimes()
}

func main() {
	StartLoadBalanceService()
	router := gin.Default()
	// 实现中间价请求
	router = Route.RouteMiddleware(router)
	// 登录注册相关的 group
	Route.UserRouter(router)
	Route.Authenciate(router)
	// 鉴权转发的代理 group
	router.Run(configs.ConfigData.LoadBalanceService.Host + configs.ConfigData.LoadBalanceService.Port)
}
