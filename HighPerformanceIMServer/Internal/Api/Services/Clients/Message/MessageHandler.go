package Message

import (
	"HighPerformanceIMServer/Configs"
	"HighPerformanceIMServer/DataModels/ApiRequests"
	"HighPerformanceIMServer/Internal/Api/Services/Dispatch"
	"HighPerformanceIMServer/Packages/Date"
	"HighPerformanceIMServer/Packages/Enums"
	"HighPerformanceIMServer/Packages/Utils"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/valyala/fastjson"
	"go.uber.org/zap"
)

// ValidationMessage 是一个方法，用于验证消息并进行处理
// 参数message是一个[]byte类型的消息
// 返回一个error类型、两个[]byte类型、一个int类型和一个string类型的结果
func (MH *Handler) ValidationMessage(message []byte) (error, []byte, []byte, int, string) {
	// 声明一个fastjson.Parser变量parser
	var parser fastjson.Parser
	if message == nil {
		// 这里把刷新网页当成一种WsPing消息
		// 如果是WsPing类型的消息，则返回ping响应消息
		zap.S().Infoln("刷新网页中！")
		return nil, []byte(`{"msg_code":1004,"message":"ping"}`),
			[]byte(``), PING, ""
	}
	// 将消息转换为字符串，并使用parser解析为fastjson.Value类型的messageJsoned
	messageJsoned, _ := parser.Parse(string(message))

	// 检查消息是否为WsPing类型的消息
	if messageCode, _ := messageJsoned.Get("msg_code").Int(); messageCode == Enums.WsPing {
		// 如果是WsPing类型的消息，则返回ping响应消息
		return nil, []byte(`{"msg_code":1004,"message":"ping"}`),
			[]byte(``), PING, ""
	}

	// 检查消息长度是否为0
	if len(message) == 0 {
		// 如果消息长度为0，则返回空消息错误响应
		return errors.New(""), []byte(`{"msg_code":1004,"message":"请勿发送空消息！"}`),
			[]byte(``), ERROR, ""
	}

	// 构建PrivateMessageRequest类型的params结构体
	params := ApiRequests.PrivateMessageRequest{
		MsgId:       Date.TimeUnixNano(),
		MsgClientId: messageJsoned.GetInt64("msg_client_id"),
		MsgCode:     Enums.WsChatMessage,
		FromID:      messageJsoned.GetInt64(),
		ToID:        messageJsoned.GetInt64(),
		MsgType:     messageJsoned.GetInt(),
		ChannelType: messageJsoned.GetInt(),
		Message:     messageJsoned.Get("msg_client_id").String(),
		SendTime:    Date.NewDate(),
		Data:        messageJsoned.Get("data").String(),
	}

	// 使用validator验证params结构体的字段是否满足规则
	if err := validator.New().Struct(params); err != nil {
		// 如果验证失败，则返回用户消息解析异常错误响应
		return err, []byte(`{"msg_code":500,"message":"用户消息解析异常"}`), []byte(``), ERROR, ""
	}

	// 构建AckMsg类型的ackMessage结构体
	ackMessage := AckMsg{
		Ack:         1,
		MsgCode:     Enums.WsAck,
		MsgId:       params.MsgId,
		MsgClientId: params.MsgClientId,
	}
	// 调用MessageHandler的GetAckMessages方法，获取ackMessage的字符串形式
	ackMessageString := MH.GetAckMessages(ackMessage)

	// 检查节点是否存在
	ok, node := IsNode(params.ToID)
	var chatMessage string
	if ok {
		// 如果节点存在，则调用MessageHandler的GetPrivateChatMessages方法，构建私聊消息字符串
		chatMessage = MH.GetPrivateChatMessages(params, true)
		return nil, []byte(chatMessage), []byte(ackMessageString),
			FORWARDING, node
	} else {
		// 如果节点不存在，则调用MessageHandler的GetPrivateChatMessages方法，构建私聊消息字符串
		chatMessage = MH.GetPrivateChatMessages(params, false)
		return nil, []byte(chatMessage), []byte(ackMessageString),
			params.ChannelType, ""
	}
}

// GetAckMessages 获取确认消息的方法
// 参数 ack 是 AckMsg 结构体，表示确认消息的参数
// 返回值 string 表示获取到的确认消息
func (MH *Handler) GetAckMessages(ack AckMsg) string {
	msg := fmt.Sprintf(`{"ack": %d,"msg_code": %d,"msg_id": %d,"msg_client_id": %d,}`, 1, ack.MsgCode, ack.MsgId, ack.MsgClientId)
	return msg
}

// IsNode IsNode是一个函数，用于检查节点是否存在
// 参数toId是一个int64类型的节点ID
// 返回一个bool类型和一个string类型的结果
func IsNode(toId int64) (bool, string) {
	// 检查配置文件中的Server的ServiceOpen字段是否为true
	if Configs.ConfigData.Server.ClusterOpen {
		// 创建一个DispatchService变量dService用于处理节点的调度
		var dService Dispatch.Service
		// 调用dService的IsDispatchNode方法，将toId转换为字符串，并将结果分配给ok和node变量
		ok, node := dService.IsDispatchNode(Utils.Int64ToString(toId))
		// 如果ok为true且node不为空字符串，则表示节点存在，返回true和节点信息
		if ok && node != "" {
			return true, node
		}
	}
	// 如果节点不存在或者为空字符串，则返回false和空字符串
	return false, ""
}

// GetPrivateChatMessages GetPrivateChatMessages是一个方法，用于获取私聊消息
// 参数message是一个PrivateMessageRequest类型的结构体，包含了私聊消息的相关信息
// 参数isGrpcMessage是一个bool类型的变量，表示是否是gRPC消息
// 返回一个string类型的结果，表示私聊消息的字符串形式
func (MH *Handler) GetPrivateChatMessages(message ApiRequests.PrivateMessageRequest, isGrpcMessage bool) string {
	// 格式化私聊消息，构建消息字符串
	msg := fmt.Sprintf(`{
		"msg_id": %d,"msg_client_id": %d,"msg_code": %d,"from_id": %d,
		"to_id": %d,"msg_type": %d,"channel_type": %d,"message": "%s","data": "%s"
	}`, message.MsgId, message.MsgClientId, message.MsgCode, message.FromID, message.ToID,
		message.MsgType, message.ChannelType, message.Message, message.Data)

	// 如果是gRPC消息，则直接返回消息字符串
	if isGrpcMessage {
		return msg
	} else {
		// 否则，构建完整的消息字符串，并返回
		msgString := fmt.Sprintf(`{
			"receive_id": "%d",
			"channel_type": %d,
			"msg": %s
		}`, message.ToID, message.ChannelType, msg)
		return msgString
	}
}
