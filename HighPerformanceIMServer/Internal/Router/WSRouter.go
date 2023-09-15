package Router

import (
	"HighPerformanceIMServer/Internal/Api/Services/WSHandler"
	"HighPerformanceIMServer/Internal/Middleware"
	"github.com/gin-gonic/gin"
)

// RegisterWSRouters 用于注册 WebSocket 路由到 Gin 引擎。
func RegisterWSRouters(router *gin.Engine) {
	WSService := new(WSHandler.WSService)            // 创建一个 WSService 实例
	ws := router.Group("/im").Use(Middleware.CORS()) // 创建一个路由组 "/im"，应用身份验证中间件和跨域中间件
	{
		ws.GET("/connect", WSService.UserConnect) // 注册 GET 请求 "/im/connect" 到 UserConnect 方法
	}
}
