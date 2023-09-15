package NSQQueue

import (
	"HighPerformanceIMServer/Configs"
	"HighPerformanceIMServer/Packages/MessageQueue"
	"HighPerformanceIMServer/Packages/MessageQueue/NSQQueue"
	"go.uber.org/zap"
)

// 定义了三个常量：ChannelOfflineTopic、ChannelGroupOfflineTopic和ChannelNodeTopic
// ChannelOfflineTopic表示离线私人消息通道
// ChannelGroupOfflineTopic表示离线群组消息通道
// ChannelNodeTopic表示节点消息通道
const (
	ChannelOfflineTopic      = "channel-offline-private"
	ChannelGroupOfflineTopic = "channel-offline-group"
	ChannelNodeTopic         = "channel-node"
)

// ConsumersPrivateMessageInit 初始化离线私人消息消费者
func ConsumersPrivateMessageInit() {
	// 调用NSQQueue的NewConsumers方法创建离线私人消息消费者
	if err := NSQQueue.NewConsumers(MessageQueue.OfflinePrivateTopic, ChannelOfflineTopic,
		Configs.ConfigData.Nsq.LookupHost); err != nil {
		// 创建消费者失败，打印错误信息
		zap.S().Infoln("new nsq consumer failed", err)
		return
	}
	// 无限循环，保持消费者运行
	select {}
}

// ConsumersGroupMessageInit 初始化离线群组消息消费者
func ConsumersGroupMessageInit() {
	// 调用NSQQueue的NewGroupConsumers方法创建离线群组消息消费者
	if err := NSQQueue.NewGroupConsumers(MessageQueue.OfflineGroupTopic, ChannelGroupOfflineTopic,
		Configs.ConfigData.Nsq.LookupHost); err != nil {
		// 创建消费者失败，打印错误信息
		zap.S().Infoln("new nsq consumer failed", err)
		return
	}
	// 无限循环，保持消费者运行
	select {}
}

// NodeInit 初始化节点消息消费者
func NodeInit() {
	// 调用NSQQueue的NewConsumers方法创建节点消息消费者
	if err := NSQQueue.NewConsumers(MessageQueue.OfflinePrivateTopic, ChannelNodeTopic,
		Configs.ConfigData.Nsq.LookupHost); err != nil {
		// 创建消费者失败，打印错误信息
		zap.S().Infoln("new nsq consumer failed", err)
		return
	}
	// 无限循环，保持消费者运行
	select {}
}
