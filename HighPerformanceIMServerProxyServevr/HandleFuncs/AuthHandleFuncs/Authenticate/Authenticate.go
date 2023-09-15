package Authenticate

import (
	"HighPerformanceIMServerProxyServevr/HandleFuncs/AUTHGRPC"
	"HighPerformanceIMServerProxyServevr/packages/Response"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Auth() gin.HandlerFunc {
	// 返回一个匿名函数作为中间件函数
	return func(ctx *gin.Context) {
		// 从请求的 URL 参数和请求头中获取 token
		token := ctx.DefaultQuery("token", ctx.GetHeader("authorization"))
		// 调用 ValidatedToken 函数验证 token 的有效性，返回验证错误和处理后的 token
		token, validatedTokenError := ValidatedToken(token)
		// 如果验证错误不为空，则表示 token 无效
		if validatedTokenError != nil {
			// 使用 zap 记录错误日志
			zap.S().Errorln(validatedTokenError)
			// 返回未授权的错误响应，并设置 HTTP 状态码为 401
			Response.ErrorResponse(http.StatusUnauthorized, validatedTokenError.Error()).SetHttpCode(http.StatusUnauthorized).WriteTo(ctx)
			// 终止后续中间件和请求处理函数的执行
			ctx.Abort()
			return
		}
		response, err := AUTHGRPC.AuthHandler(token)
		if err != nil {
			// 使用 zap 记录错误日志
			zap.S().Errorln(err)
			// 返回未授权的错误响应，并设置 HTTP 状态码为 401
			Response.ErrorResponse(http.StatusUnauthorized, err.Error()).SetHttpCode(http.StatusUnauthorized).WriteTo(ctx)
			// 终止后续中间件和请求处理函数的执行
			ctx.Abort()
			return
		}
		zap.S().Infoln("成功了")
		// 将 claims 中的 ID、UID、Name 设置到 gin 的上下文中，供后续的处理函数使用
		ctx.Set("id", response.ID)
		ctx.Set("uid", response.UID)
		ctx.Set("name", response.Name)
		ctx.Next()
	}
}

func ValidatedToken(token string) (string, error) {
	var err error // 声明一个error类型的变量err

	if len(token) == 0 {
		err = errors.New("Token不能为空")
		return err.Error(), err // 如果token为空，则返回err和"Token 不能为空"的错误信息
	}

	t := strings.Split(token, "Bearer ") // 使用空格分割token字符串，并将结果赋值给t变量
	if len(t) > 1 {
		return t[1], nil // 如果t的长度大于1，则返回nil和t的第二个元素
	} else {
		return token, nil // 否则，返回nil和原始的token值
	}
}
