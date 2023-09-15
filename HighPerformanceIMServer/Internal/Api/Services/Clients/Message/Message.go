package Message

import "HighPerformanceIMServer/DataModels/ApiRequests"

// Client 结构体用于表示消息客户端的请求参数
type Client struct {
	ReceiveId   string  `json:"receive_id"`   // 接收消息的用户ID
	ChannelType int     `json:"channel_type"` // 频道类型
	Msg         Message `json:"msg"`          // 消息内容
}

// Message 用户发送的消息
type Message struct {
	MsgId       int64       `json:"msg_id"`        // 服务端消息唯一id
	MsgClientId int64       `json:"msg_client_id"` // 客户端消息唯一id
	MsgCode     int         `json:"msg_code"`      // 定义的消息code
	FromID      int64       `json:"from_id"`       // 发消息的人
	ToID        int64       `json:"to_id"`         // 接收消息人的id
	Uid         string      `json:"uid"`           // uid
	ToUid       string      `json:"to_uid"`        // to uid
	MsgType     int         `json:"msg_type"`      // 消息类型 例如 1.文本 2.语音 3.文件
	ChannelType int         `json:"channel_type"`  // 频道类型 1.私聊 2.频道 3.广播
	Message     string      `json:"message"`       // 消息
	SendTime    string      `json:"send_time"`     // 消息发送时间
	Data        interface{} `json:"data"`          // 自定义携带的数据
}

// AckMsg ack消息
type AckMsg struct {
	Ack         int   `json:"ack"`           // 1.消息已经投递到服务器了
	MsgCode     int   `json:"msg_code"`      // 1.消息已经投递到服务器了
	MsgId       int64 `json:"msg_id"`        //服务器生成的消息id
	MsgClientId int64 `json:"msg_client_id"` //客户端生成的消息id
}

// CreateFriendMessage 私聊内容
type CreateFriendMessage struct {
	MsgCode     int    `json:"msg_code"`    // 定义的消息code
	ID          int64  `json:"id"`          // 定义的消息code
	FromId      int64  `json:"from_id"`     // 发消息的人
	Status      int    `json:"status"`      // 发消息的人
	CreatedAt   string `json:"created_at"`  // 发消息的人
	ToID        int64  `json:"to_id"`       // 接收消息人的id
	Information string `json:"information"` //内容
	Users       Users  `json:"users"`       //请求人
}

// Users 用户信息
type Users struct {
	Name   string `json:"name"`   // 用户名
	ID     int64  `json:"id"`     // 用户ID
	Avatar string `json:"avatar"` // 用户头像
}

// PingMessage 心跳消息
type PingMessage struct {
	MsgCode int    `json:"msg_code"` // 消息代码
	Message string `json:"message"`  // 消息内容
}

// BroadcastMessages TODO
type BroadcastMessages struct {
}

// Interface MessageInterface 是一个消息接口
type Interface interface {
	// ValidationMsg 验证消息的方法
	// 参数 msg 是待验证的消息的字节流
	// 返回值 error 表示验证过程中的错误，string 表示验证结果
	ValidationMsg(msg []byte) (error, string)

	// GetPrivateChatMessages 获取私聊消息的方法
	// 参数 message 是 PrivateMessageRequest 结构体，表示私聊消息的请求参数
	// 参数 isGrpcMessage 表示是否为 gRPC 消息
	// 返回值 string 表示获取到的私聊消息
	GetPrivateChatMessages(message ApiRequests.PrivateMessageRequest, isGrpcMessage bool) string

	// GetAckMessages 获取确认消息的方法
	// 参数 ack 是 AckMsg 结构体，表示确认消息的参数
	// 返回值 string 表示获取到的确认消息
	GetAckMessages(ack AckMsg) string
}
type Handler struct {
}

const (
	ERROR      = 0 // 错误消息类型
	PRIVATE    = 1 // 私聊消息类型
	GROUP      = 2 // 群聊消息类型
	PING       = 3 // 心跳消息类型
	FORWARDING = 4 // 转发消息类型
)
