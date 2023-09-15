package MessageGRPC

import (
	"HighPerformanceIMServer/DataModels/ApiRequests"
	"HighPerformanceIMServer/Packages/Date"
	"HighPerformanceIMServer/Packages/Enums"
	"context"
	"fmt"
	"go.uber.org/zap"
)

type IMGRPCMessage struct {
}

func (IMGRPCMessage) mustEmbedUnimplementedImMessageServer() {}

// SendMessageHandler 是 IMGroupMessage 结构体的方法，用于处理发送消息的请求并返回响应。
func (IMGRPCMessage) SendMessageHandler(ctx context.Context, request *SendMessageRequest) (*SendMessageResponse, error) {
	zap.S().Errorln(request.Message) // 输出请求中的消息内容

	params := ApiRequests.PrivateMessageRequest{
		MsgId:       Date.TimeUnixNano(),      // 生成消息ID
		MsgCode:     Enums.WsChatMessage,      // 消息代码
		MsgClientId: request.MsgClientId,      // 消息客户端ID
		FromID:      request.FromId,           // 发送方ID
		ToID:        request.ToId,             // 接收方ID
		ChannelType: int(request.ChannelType), // 渠道类型
		MsgType:     int(request.MsgType),     // 消息类型
		Message:     request.Message,          // 消息内容
		SendTime:    Date.NewDate(),           // 发送时间
		Data:        request.Data,             // 数据
	}

	messageString := GetGrpcPrivateChatMessages(params) // 获取 gRPC 私聊消息的字符串表示
	switch request.ChannelType {
	case 1:
		fmt.Println(messageString) // 打印消息字符串
	case 2:
		// TODO: 添加对其他渠道类型的处理
	}

	return &SendMessageResponse{
		Code:    200,       // 响应代码
		Message: "Success", // 响应消息
	}, nil
}

// GetGrpcPrivateChatMessages 用于获取 gRPC 私聊消息的字符串表示。
func GetGrpcPrivateChatMessages(message ApiRequests.PrivateMessageRequest) string {
	msg := fmt.Sprintf(`{
		"msg_id": %d,
		"msg_client_id": %d,
		"msg_code": %d,
		"form_id": %d,
		"to_id": %d,
		"msg_type": %d,
		"channel_type": %d,
		"message": %s,
		"data": %s
	}`, message.MsgId, message.MsgClientId, message.MsgCode, message.FromID, message.ToID, message.MsgType, message.ChannelType, message.Message, message.Data)

	return msg
}
