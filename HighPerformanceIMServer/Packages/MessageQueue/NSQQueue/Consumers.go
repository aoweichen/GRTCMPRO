package NSQQueue

import (
	"HighPerformanceIMServer/Packages/MessageOfflineDao"
	"github.com/nsqio/go-nsq"
	"go.uber.org/zap"
)

var (
	OfflineMessageSave = new(MessageOfflineDao.OfflineMessageDao)
)

// PrivateHandler 定义了一个名为PrivateHandler的结构体
type PrivateHandler struct {
}

// GroupHandler 定义了一个名为GroupHandler的结构体
type GroupHandler struct {
}

// HandleMessage 定义了PrivateHandler结构体的HandleMessage方法
// 该方法用于处理私聊消息
// 方法接收一个messageNSQ参数，表示接收到的NSQ消息
// 方法返回一个error类型的值，表示处理过程中是否出现错误
func (PH *PrivateHandler) HandleMessage(messageNSQ *nsq.Message) error {
	// 调用OfflineMessageSave的PrivateOfflineMessageSave方法，将消息内容保存为私聊离线消息
	OfflineMessageSave.PrivateOfflineMessageSave(string(messageNSQ.Body))
	return nil
}

// HandleMessage 定义了GroupHandler结构体的HandleMessage方法
// 该方法用于处理群聊消息
// 方法接收一个messageNSQ参数，表示接收到的NSQ消息
// 方法返回一个error类型的值，表示处理过程中是否出现错误
func (GH *GroupHandler) HandleMessage(messageNSQ *nsq.Message) error {
	// 调用OfflineMessageSave的GroupOfflineMessageSave方法，将消息内容保存为群聊离线消息
	OfflineMessageSave.GroupOfflineMessageSave(string(messageNSQ.Body))
	return nil
}

// NewConsumers 定义了一个名为NewConsumers的方法
// 该方法用于创建NSQ消费者，并连接到指定的NSQLookupd地址
// 方法接收三个参数：topic表示消费的主题，channel表示消费的通道，addr表示NSQLookupd地址
// 方法返回一个error类型的值，表示创建消费者过程中是否出现错误
func NewConsumers(topic string, channel string, addr string) error {
	// 创建NSQ配置对象
	configNSQ := nsq.NewConfig()

	// 创建NSQ消费者
	nsqConsumer, err := nsq.NewConsumer(topic, channel, configNSQ)
	if err != nil {
		// 创建消费者失败，打印错误信息并返回错误
		zap.S().Errorln("create consumer failed err ", err)
		return err
	}

	// 创建私聊消息处理器
	consumer := &PrivateHandler{}

	// 将私聊消息处理器添加到NSQ消费者中
	nsqConsumer.AddHandler(consumer)

	// 连接到指定的NSQLookupd地址
	if err := nsqConsumer.ConnectToNSQLookupd(addr); err != nil {
		// 连接失败，打印错误信息并返回错误
		zap.S().Errorln(err)
		return err
	} else {
		// 连接成功，返回nil表示没有错误
		return nil
	}
}

// NewGroupConsumers 定义了一个名为NewGroupConsumers的方法
// 该方法用于创建NSQ群组消费者，并连接到指定的NSQLookupd地址
// 方法接收三个参数：topic表示消费的主题，channel表示消费的通道，addr表示NSQLookupd地址
// 方法返回一个error类型的值，表示创建群组消费者过程中是否出现错误
func NewGroupConsumers(topic string, channel string, addr string) error {
	// 创建NSQ配置对象
	configNSQ := nsq.NewConfig()

	// 创建NSQ群组消费者
	nsqConsumer, err := nsq.NewConsumer(topic, channel, configNSQ)
	if err != nil {
		// 创建消费者失败，打印错误信息并返回错误
		zap.S().Errorln("create consumer failed err ", err)
		return err
	}

	// 创建群组消息处理器
	consumer := &GroupHandler{}

	// 将群组消息处理器添加到NSQ消费者中
	nsqConsumer.AddHandler(consumer)

	// 连接到指定的NSQLookupd地址
	if err := nsqConsumer.ConnectToNSQLookupd(addr); err != nil {
		// 连接失败，打印错误信息并返回错误
		zap.S().Errorln(err)
		return err
	} else {
		// 连接成功，返回nil表示没有错误
		return nil
	}
}
