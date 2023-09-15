package Middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"runtime/debug"
)

// Recover Recover函数是一个中间件函数，用于捕获处理请求过程中的任何panic或recover错误
func Recover(ctx *gin.Context) {
	// 定义一个defer函数，确保无论是否发生错误，都会执行这个函数
	defer func() {
		// 如果发生错误，调用recover函数进行异常恢复
		if r := recover(); r != nil {
			// 将异常转换为字符串，以便记录和显示S
			var errorString = ErrorToString(r)
			// 将错误信息记录到日志
			zap.S().Errorln(errorString)
			// 打印异常的堆栈跟踪信息
			debug.PrintStack()

			// 返回一个HTTP 500（内部服务器错误）响应，并将错误信息作为响应体
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": errorString,
			})
			// 中断当前的请求处理流程
			ctx.Abort()
		}
	}()
	// 继续处理请求
	ctx.Next()
}

// ErrorToString errorToString函数是一个辅助函数，用于将输入的错误转换为字符串
func ErrorToString(r interface{}) string {
	// 根据输入的类型，将其转换为字符串
	switch v := r.(type) {
	case error:
		// 如果输入是一个error类型，调用其Error()方法获取错误信息
		return v.Error()
	default:
		// 否则，直接将输入转换为字符串
		return r.(string)
	}
}
