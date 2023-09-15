package Manager

import (
	"HighPerformanceIMServer/Internal/Api/Services/Clients/GRPCMessage"
	"HighPerformanceIMServer/Internal/Api/Services/Clients/Message"
	"HighPerformanceIMServer/Internal/Api/Services/Dispatch"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"sync"
)

// IMClient IM客户端
type IMClient struct {
	ClientID string          // 客户端用户id
	UUID     string          // 用户唯一id
	Socket   *websocket.Conn // 当前socket握手对象
	Send     chan []byte     // 当前用户发送的消息
	Mux      sync.RWMutex    // 互斥锁
	Identity int             // 身份 1.游客 2.用户
}

type WSClientInterface interface {
	// WSReadMessage websocket 读消息
	WSReadMessage()
	// WSWriteMessage websocket 发消息
	WSWriteMessage()
	// WSClose 关闭WS连接
	WSClose()
}

var (
	messageHandler    Message.Handler
	gRPCMessageClient GRPCMessage.MessageService
	dispatchNode      Dispatch.Service
)

// NewWSClient 是一个构造函数，用于创建一个新的 IMClient 实例，并返回指针类型的 IMClient
// 参数 ID 是客户端的唯一标识符，uid 是客户端的用户标识符，identity 是客户端的身份信息，conn 是 websocket 连接对象
func NewWSClient(ID string, uid string, identity int, conn *websocket.Conn) *IMClient {
	// 创建一个新的 IMClient 实例，并赋值给变量 client
	client := &IMClient{
		ClientID: ID,
		UUID:     uid,
		Socket:   conn,
		Send:     make(chan []byte),
		Identity: identity,
	}
	// 返回指向 client 的指针
	return client
}

// WSReadMessage 是 IMClient 类型的方法，用于读取来自客户端的消息并进行处理
func (IMC *IMClient) WSReadMessage() {
	// 在方法结束时，将当前客户端从 IMMessageClientManager 的 CloseClients 通道中发送出去，并关闭客户端的连接
	defer func() {
		IMMessageClientManager.CloseClients <- IMC
		_ = IMC.Socket.Close()
	}()
	for {
		// 从客户端的连接中读取消息
		_, message, err := IMC.Socket.ReadMessage()
		if err != nil {
			// 如果读取消息出错，则将当前客户端从 IMMessageClientManager 的 CloseClients 通道中发送出去，并关闭客户端的连接
			IMMessageClientManager.CloseClients <- IMC
			IMC.WSClose()
			// 客户端移除后直接return，以防进入后面的逻辑
			break
		}

		// 验证并处理消息
		err, messageBytes, ackMessage, channel, node := messageHandler.ValidationMessage(message)
		if err != nil {
			// 如果消息验证出错，则将消息的字节表示打印出来，并将消息通过客户端的连接发送回去
			zap.S().Infoln(string(messageBytes))
			_ = IMC.Socket.WriteMessage(websocket.TextMessage, messageBytes)
		} else {
			switch channel {
			case Message.PRIVATE:
				// 如果消息类型是私聊，则将确认消息通过客户端的连接发送回去，并将消息的字节表示发送到 IMMessageClientManager 的 PrivateChannel 通道中
				_ = IMC.Socket.WriteMessage(websocket.TextMessage, ackMessage)
				IMMessageClientManager.PrivateChannel <- messageBytes
			case Message.GROUP:
				// 如果消息类型是群聊，则将确认消息通过客户端的连接发送回去
				_ = IMC.Socket.WriteMessage(websocket.TextMessage, ackMessage)
			case Message.PING:
				// 如果消息类型是PING，则将确认消息通过客户端的连接发送回去
				_ = IMC.Socket.WriteMessage(websocket.TextMessage, ackMessage)
			case Message.FORWARDING:
				// 如果消息类型是转发，则将确认消息通过客户端的连接发送回去，并调用 gRPCMessageClient 的 SendGRPCMessage 方法发送 gRPC 消息
				_ = IMC.Socket.WriteMessage(websocket.TextMessage, ackMessage)
				gRPCMessageClient.SendGRPCMessage(string(messageBytes), node)
			default:
				// 如果消息类型是其他类型，则将确认消息通过客户端的连接发送回去，并将消息的字节表示发送到 IMMessageClientManager 的 BroadcastChannel 通道中
				_ = IMC.Socket.WriteMessage(websocket.TextMessage, ackMessage)
				IMMessageClientManager.BroadcastChannel <- messageBytes
			}
		}
	}
}

// WSWriteMessage 是 IMClient 类型的方法，用于向客户端写入消息
func (IMC *IMClient) WSWriteMessage() {
	// 在方法结束时，调用 WSClose 方法关闭客户端的连接
	defer IMC.WSClose()

	for {
		select {
		case message, ok := <-IMC.Send:
			// 如果通道关闭，则向客户端发送关闭消息，然后返回
			if !ok {
				_ = IMC.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			// 向客户端发送文本消息
			_ = IMC.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

// WSClose 是 IMClient 类型的方法，用于关闭客户端的连接
func (IMC *IMClient) WSClose() {
	// 调用 dispatchNode 包中的 DeleteDispatchNode 方法，从 dispatchNode 中删除 IMC 的 ClientID
	dispatchNode.DeleteDispatchNode(IMC.ClientID)
	// 关闭客户端的连接
	_ = IMC.Socket.Close()
}
