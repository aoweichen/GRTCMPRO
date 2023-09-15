package MessageQueue

// OfflinePrivateTopic 定义了全局变量OfflinePrivateTopic，表示离线私人消息的主题
var OfflinePrivateTopic = "offline_private_message" //离线私人消息频道

// OfflineGroupTopic 定义了全局变量OfflineGroupTopic，表示离线群组消息的主题
var OfflineGroupTopic = "offline_group_message" //离线私人消息频道

// MessageProducerQueueInterface 定义了MessageProducerQueueInterface接口
// 该接口用于发送消息
type MessageProducerQueueInterface interface {
	SendMessage(msg []byte)
}
