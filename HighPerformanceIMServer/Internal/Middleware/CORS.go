package Middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CORS() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		reqMethod := ginCtx.Request.Method

		ginCtx.Header("Access-Control-Allow-Origin", "*")
		// 设置允许的请求头字段
		ginCtx.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token, x-token")
		// 设置了允许的请求方法
		ginCtx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
		// 设置暴露给客户端的响应头字段
		ginCtx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		// 允许携带凭证信息
		ginCtx.Header("Access-Control-Allow-Credentials", "true")
		ginCtx.Header("X-Content-Type-Options", "nosniff")
		// 设置允许的来源,设置为"*"表示允许任何来源的请求。

		// 如果请求方法为OPTIONS，则返回204 No Content状态码

		if reqMethod == "OPTIONS" {
			ginCtx.AbortWithStatus(http.StatusNoContent)
		}
	}
}
