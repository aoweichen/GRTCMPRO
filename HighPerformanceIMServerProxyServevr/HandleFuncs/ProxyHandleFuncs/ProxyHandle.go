package ProxyHandleFuncs

import (
	"HighPerformanceIMServerProxyServevr/packages/LoadBalanceHTTP"
	"fmt"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 创建自定义的 RoundTripper

// HandleProxy 处理代理请求的函数
func HandleProxy(ctx *gin.Context) {
	query := url.Values{}
	query.Set("id", fmt.Sprintf("%v", ctx.MustGet("id")))
	query.Set("uid", fmt.Sprintf("%v", ctx.MustGet("uid")))
	query.Set("name", fmt.Sprintf("%v", ctx.MustGet("name")))

	zap.S().Infoln("进入处理代理请求的函数：HandleProxy")
	// 获取代理请求的代理路径
	proxyPath := ctx.Param("proxyPath")
	// 从上下文中获取代理路径参数，用于后续处理代理路径
	zap.S().Infoln("代理路径参数 proxyPath: ", proxyPath)
	// 获取代理请求的URL
	proxyURL := ctx.Request.URL.Path
	// 获取代理请求的URL路径，用于后续处理
	zap.S().Infoln("代理请求的URL => proxyURL：", proxyURL)
	// 从上下文中获取目标URL
	targetURL := ctx.MustGet("targetURL").(*url.URL)
	// 从上下文中获取目标URL，这个URL是请求要代理到的目标地址
	zap.S().Infoln("targetURL: ", targetURL)
	// 创建自定义的 RoundTripper
	// 创建反向代理
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	// 使用目标URL创建一个反向代理对象

	// 设置请求的Host、Path、RawPath和RawQuery
	ctx.Request.Host = targetURL.Host
	// 设置请求的Host字段为目标URL的Host字段，用于正确路由代理请求
	ctx.Request.URL.Path = proxyPath
	zap.S().Infoln("proxyPath", proxyPath)
	// 设置请求的URL路径为代理路径，用于正确路由代理请求
	ctx.Request.URL.RawPath = query.Encode()
	// 清空RawPath字段，确保URL的路径信息正确
	ctx.Request.URL.RawQuery = strings.Replace(ctx.Request.URL.RawQuery, proxyURL[len(proxyPath):], "", 1)
	// 从请求的URL中去掉代理路径部分，确保查询参数正确

	// 设置X-Forwarded-For请求头
	clientIP := ctx.ClientIP()
	// 获取客户端的IP地址
	if clientIP != "" {
		ctx.Request.Header.Set("X-Forwarded-For", clientIP)
		// 设置X-Forwarded-For请求头，用于标识真实客户端的IP地址
	}
	start := time.Now()
	// 执行反向代理
	proxy.ServeHTTP(ctx.Writer, ctx.Request)
	duration := time.Since(start).Nanoseconds()
	// 使用反向代理处理请求，将代理请求转发到目标URL，并将响应写入响应写入器

	// 更新响应时间
	// 更新目标URL的响应时间
	LoadBalanceHTTP.UpdateResponseTime(targetURL.String(), duration)
}
