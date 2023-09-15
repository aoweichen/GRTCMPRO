package NSQQueue

import (
	"HighPerformanceIMServer/Packages/MessageQueue"
	"HighPerformanceIMServer/Packages/MessageQueue/NSQQueue"
	"go.uber.org/zap"
)

// ProducerQueue 定义了全局变量ProducerQueue
var ProducerQueue MessageProducerQueue

// MessageProducerQueue 定义了MessageProducerQueue结构体
type MessageProducerQueue struct {
}

// SendMessage 实现了MessageProducerQueueInterface接口的SendMessage方法
// 该方法用于发送离线私人消息
func (MPQ *MessageProducerQueue) SendMessage(message []byte) {
	// 调用NSQQueue的PublishMessage方法发送消息到离线私人消息主题
	if err := NSQQueue.PublishMessage(MessageQueue.OfflinePrivateTopic, message); err != nil {
		// 发送消息失败，打印错误信息
		zap.S().Infoln(err)
	}
}

// SendGroupMessage 实现了MessageProducerQueueInterface接口的SendGroupMessage方法
// 该方法用于发送离线群组消息
func (MPQ *MessageProducerQueue) SendGroupMessage(message []byte) {
	// 调用NSQQueue的PublishMessage方法发送消息到离线群组消息主题
	if err := NSQQueue.PublishMessage(MessageQueue.OfflineGroupTopic, message); err != nil {
		// 发送消息失败，打印错误信息
		zap.S().Infoln(err)
	}
}
