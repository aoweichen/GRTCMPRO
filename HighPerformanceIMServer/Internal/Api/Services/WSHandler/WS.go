package WSHandler

import (
	"HighPerformanceIMServer/Internal/Api/Services/Clients/Manager"
	"HighPerformanceIMServer/Internal/Api/Services/Dispatch"
	"HighPerformanceIMServer/Packages/Utils"
	"HighPerformanceIMServer/Packages/WebSocketServer"
	"github.com/gin-gonic/gin"
	"net/http"
)

type WSService struct {
}

const (
	touristsRole = 1
	userRole     = 2
)

// UserConnect 是 WSService 的方法，用于处理用户的 WebSocket 连接。
func (*WSService) UserConnect(ctx *gin.Context) {
	conn, err := WebSocketServer.App(ctx.Writer, ctx.Request) // 将 HTTP 连接升级为 WebSocket 连接
	if err != nil {
		http.Error(ctx.Writer, ctx.Errors.String(), http.StatusInternalServerError) // 如果发生错误，返回 500 错误状态码，并输出错误信息
	}

	var dispatchService Dispatch.Service                                                              // 创建 Service 实例
	id, uid := Utils.InterfaceToInt64(ctx.MustGet("id")), Utils.InterfaceToString(ctx.MustGet("uid")) // 从上下文中获取 id 和 uid

	dispatchService.SetDispatchNode(Utils.Int64ToString(id)) // 设置调度节点

	wsClient := Manager.NewWSClient(Utils.Int64ToString(id), uid, userRole, conn) // 创建一个新的 WebSocket 客户端实例，使用 id 和 uid
	Manager.IMMessageClientManager.Register <- wsClient                           // 注册 WebSocket 客户端到客户端管理器

	go wsClient.WSReadMessage()  // 启动一个协程用于读取 WebSocket 消息
	go wsClient.WSWriteMessage() // 启动一个协程用于写入 WebSocket 消息
}

// TouristsConnect 是 WSService 的方法，用于处理游客的 WebSocket 连接。
func (*WSService) TouristsConnect(ctx *gin.Context) {
	conn, err := WebSocketServer.App(ctx.Writer, ctx.Request) // 将 HTTP 连接升级为 WebSocket 连接
	if err != nil {
		http.Error(ctx.Writer, ctx.Errors.String(), http.StatusInternalServerError) // 如果发生错误，返回 500 错误状态码，并输出错误信息
	}

	id := ctx.Query("token_id") // 从请求参数中获取 token_id，用于生成客户端的唯一标识符（生成规则 - ip+时间）

	wsClient := Manager.NewWSClient(id, "", touristsRole, conn) // 创建一个新的 WebSocket 客户端实例
	Manager.IMMessageClientManager.Register <- wsClient         // 注册 WebSocket 客户端到客户端管理器

	go wsClient.WSReadMessage()  // 启动一个协程用于读取 WebSocket 消息
	go wsClient.WSWriteMessage() // 启动一个协程用于写入 WebSocket 消息
}
