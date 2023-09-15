// Package WebSocketServer 实现了一个 WebSocket 服务器
package WebSocketServer

import (
	"errors"
	"github.com/gorilla/websocket"
	"net/http"
)

// 创建一个 websocket.Upgrader 对象，用于升级 HTTP 连接为 WebSocket 连接
var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024, // 设置读缓冲区大小
	WriteBufferSize: 1024, // 设置写缓冲区大小
	CheckOrigin: func(request *http.Request) bool {
		return true // 允许所有来源的请求
	},
}

// WebSocketConn 定义一个全局变量 WebSocketConn，用于保存 WebSocket 连接
var (
	WebSocketConn *websocket.Conn
)

// App 是处理 WebSocket 连接的函数，它将升级 HTTP 连接为 WebSocket 连接，并返回 WebSocket 连接对象和可能的错误
func App(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	err := errors.New("")                           // 创建一个空错误对象
	WebSocketConn, err = upgrade.Upgrade(w, r, nil) // 使用 upgrade 对象将 HTTP 连接升级为 WebSocket 连接
	return WebSocketConn, err                       // 返回 WebSocket 连接对象和错误
}
