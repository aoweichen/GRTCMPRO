package GRPCMessage

import (
	"HighPerformanceIMServer/Internal/Api/Services/Clients/GRPCMessage/GRPC/MessageGRPC"
	"HighPerformanceIMServer/Packages/Date"
	"HighPerformanceIMServer/Packages/Enums"
	"context"
	"crypto/tls"
	"github.com/valyala/fastjson"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type ServiceInterface interface {
	// SendGRPCMessage 消息发送到指定节点
	SendGRPCMessage(message []byte, node string)
}

type MessageService struct {
}

// SendGRPCMessage 发送grpcMessage服务的message消息
func (GMS *MessageService) SendGRPCMessage(message string, node string) {
	// 使用credentials.NewTLS函数创建一个基于TLS的凭证，其中的InsecureSkipVerify设置为true表示跳过对服务器证书的验证。
	insecureCreds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})

	// 创建gRPC拨号选项,不经过TLS验证
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecureCreds),
	}

	// 拨号创建与gRPC服务器的连接
	conn, err := grpc.Dial(node, opts...)
	if err != nil {
		zap.S().Errorln("grpc 服务连接失败:", err)
		panic(err)
	}

	// 延迟关闭连接
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			zap.S().Errorln("grpc conn close error:", err)
			panic(err)
		}
	}(conn)

	// 创建ImMessage客户端
	IMGRPCServiceClient := MessageGRPC.NewImMessageClient(conn)

	// 使用fastjson.Parser解析消息字符串为JSON对象
	var parser fastjson.Parser
	messageFastJson, _ := parser.Parse(message)

	// 构造发送消息的请求参数
	params := &MessageGRPC.SendMessageRequest{
		MsgId:       Date.TimeUnixNano(),
		MsgClientId: Date.TimeUnix(),
		MsgCode:     Enums.WsChatMessage,
		FromId:      messageFastJson.GetInt64("from_id"),
		ToId:        messageFastJson.GetInt64("to_id"),
		MsgType:     messageFastJson.GetInt64("msg_type"),
		ChannelType: messageFastJson.GetInt64("channel_type"),
		Message:     messageFastJson.Get("message").String(),
		SendTime:    Date.TimeUnixNano(),
		Data:        messageFastJson.Get("data").String(),
	}

	// 调用发送消息的RPC方法
	response, err := IMGRPCServiceClient.SendMessageHandler(context.Background(), params)
	if err != nil {
		zap.S().Errorln("发送 grpc 请求失败（SendMessageHandler）:", err)
		panic(err)
	}

	// 打印响应消息
	zap.S().Infoln(response.Message)
	return
}
