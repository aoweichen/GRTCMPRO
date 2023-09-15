package Route

import (
	"HighPerformanceIMServerProxyServevr/HandleFuncs/AuthHandleFuncs/Auth"
	"HighPerformanceIMServerProxyServevr/HandleFuncs/AuthHandleFuncs/Authenticate"
	"HighPerformanceIMServerProxyServevr/HandleFuncs/ProxyHandleFuncs"
	"HighPerformanceIMServerProxyServevr/Middleware"
	"HighPerformanceIMServerProxyServevr/packages/LoadBalanceHTTP"

	"github.com/gin-gonic/gin"
)

func RouteMiddleware(router *gin.Engine) *gin.Engine {
	// 解决 CORS 问题
	router.Use(Middleware.CORS())
	// 从panic中恢复
	router.Use(Middleware.Recover)

	return router
}

func UserRouter(router *gin.Engine) {
	UserGroup := router.Group("/user")
	{
		UserGroup.POST("/login", Auth.Login)
		UserGroup.POST("/register", Auth.Register)
		UserGroup.POST("/sendEmailCode", Auth.SendeEmail)
	}

}

func Authenciate(router *gin.Engine) {
	router.Use(Authenticate.Auth())
	router.Use(LoadBalanceHTTP.GetTargetIMServerURLFromConsulMiddleware)
	AuthGroup := router.Group("/auth")
	{
		AuthGroup.Any("/*proxyPath", ProxyHandleFuncs.HandleProxy)
	}
}
