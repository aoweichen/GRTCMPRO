package MessageOfflineDao

import (
	"HighPerformanceIMServer/DataModels/Models/IMOfflineMessage"
	"HighPerformanceIMServer/Internal/DAO/MYSQL"
	"HighPerformanceIMServer/Packages/Date"
	"github.com/valyala/fastjson"
	"sync"
)

// DataInterface 定义了一个名为DataInterface的接口
type DataInterface interface {
	PrivateOfflineMessageSave(string)                                         // 保存私聊离线消息的方法，接收一个string类型参数
	GroupOfflineMessageSave(string)                                           // 保存群聊离线消息的方法，接收一个string类型参数
	PullPrivateOfflineMessage(string) []IMOfflineMessage.ImOfflineMessages    // 拉取私聊离线消息的方法，接收一个int64类型参数，返回一个IMOfflineMessage.ImOfflineMessages切片
	PullGroupOfflineMessage(string) []IMOfflineMessage.ImGroupOfflineMessages // 拉取群聊离线消息的方法，接收一个int64类型参数，返回一个IMOfflineMessage.ImOfflineMessages切片
	UpdateOfflineMessageStatus(string, int)                                   // 更新私聊离线消息状态的方法，接收一个int64类型参数
	PrivateOfflineMessageSaveOptimization(messages []string)                  // 保存私聊离线消息的方法，接收一个string类型参数(优化合并成多条进行插入数据库 主从)
	GroupOfflineMessageSaveOptimization(messages []string)                    // 保存群聊离线消息的方法，接收一个string类型参数(优化合并成多条进行插入数据库 主从)
}

var (
	OfflineMessage OfflineMessageDao
)

type OfflineMessageDao struct {
}

// PrivateOfflineMessageSave
// 定义了一个名为PrivateOfflineMessageSave的方法，属于OfflineMessageDao类型的接收者
// 该方法用于保存私聊离线消息
// 方法接收一个message参数，表示要保存的消息内容
func (OMD *OfflineMessageDao) PrivateOfflineMessageSave(message string) {
	// 使用fastjson.Parser解析消息内容
	var parserFastJson fastjson.Parser
	messageParsed, _ := parserFastJson.Parse(message)

	// 从解析后的消息中获取接收消息的用户ID
	receivedID := messageParsed.GetInt64("to_id")

	// 创建一个IMOfflineMessage.ImOfflineMessages对象，并保存到数据库表im_offline_messages中
	MYSQL.DataBase.Table("im_offline_messages").Create(&IMOfflineMessage.ImOfflineMessages{
		Status:    0,
		SendTime:  int(Date.TimeUnix()),
		ReceiveId: receivedID,
		Message:   message,
	})
}

// PrivateOfflineMessageSaveOptimization 定义了一个名为PrivateOfflineMessageSaveOptimization的方法，属于OfflineMessageDao类型的接收者
// 该方法用于优化保存多条私聊离线消息的操作
// 方法接收一个messages参数，表示要保存的多条消息内容组成的切片
func (OMD *OfflineMessageDao) PrivateOfflineMessageSaveOptimization(messages []string) {
	// 创建一个互斥锁，用于保证并发安全
	var mutex sync.Mutex

	// 创建一个切片，用于存储离线消息对象
	var offlineMessages []*IMOfflineMessage.ImOfflineMessages

	// 遍历消息切片
	for _, message := range messages {
		// 使用fastjson.Parser解析消息内容
		var parserFastJson fastjson.Parser
		messageParsed, _ := parserFastJson.Parse(message)

		// 从解析后的消息中获取接收消息的用户ID
		receivedID := messageParsed.GetInt64("to_id")

		// 创建一个IMOfflineMessage.ImOfflineMessages对象，并添加到离线消息切片中
		offlineMessage := &IMOfflineMessage.ImOfflineMessages{
			Status:    0,
			SendTime:  int(Date.TimeUnix()),
			ReceiveId: receivedID,
			Message:   message,
		}

		offlineMessages = append(offlineMessages, offlineMessage)
	}

	// 使用互斥锁保证并发安全
	mutex.Lock()
	// 将离线消息切片保存到数据库表im_offline_messages中
	MYSQL.DataBase.Table("im_offline_messages").Create(&offlineMessages)
	mutex.Unlock()
}

// GroupOfflineMessageSave
// 定义了一个名为GroupOfflineMessageSave的方法，属于OfflineMessageDao类型的接收者
// 该方法用于保存群聊离线消息
// 方法接收一个message参数，表示要保存的消息内容
func (OMD *OfflineMessageDao) GroupOfflineMessageSave(message string) {
	// 使用fastjson.Parser解析消息内容
	var parserFastJson fastjson.Parser
	messageParsed, _ := parserFastJson.Parse(message)

	// 从解析后的消息中获取接收消息的用户ID
	userID := messageParsed.GetInt("user_id")

	// 创建一个IMOfflineMessage.ImGroupOfflineMessages对象，并保存到数据库表im_offline_messages中
	MYSQL.DataBase.Table("im_offline_messages").Create(&IMOfflineMessage.ImGroupOfflineMessages{
		Status:    0,
		SendTime:  int(Date.TimeUnix()),
		ReceiveId: userID,
		Message:   message,
	})
}

// GroupOfflineMessageSaveOptimization 定义了一个名为GroupOfflineMessageSave的方法，属于OfflineMessageDao类型的接收者
// 该方法用于保存多条群聊离线消息
// 方法接收一个messages参数，表示要保存的多条消息内容组成的切片
func (OMD *OfflineMessageDao) GroupOfflineMessageSaveOptimization(messages []string) {
	// 创建一个互斥锁，用于保证并发安全
	var mutex sync.Mutex
	// 创建一个切片，用于存储群聊离线消息对象
	var offlineMessages []*IMOfflineMessage.ImGroupOfflineMessages

	// 遍历消息切片
	for _, message := range messages {
		// 使用fastjson.Parser解析消息内容
		var parserFastJson fastjson.Parser
		messageParsed, _ := parserFastJson.Parse(message)

		// 从解析后的消息中获取接收消息的用户ID
		userID := messageParsed.GetInt("user_id")

		// 创建一个IMOfflineMessage.ImGroupOfflineMessages对象，并添加到离线消息切片中
		offlineMessage := &IMOfflineMessage.ImGroupOfflineMessages{
			Status:    0,
			SendTime:  int(Date.TimeUnix()),
			ReceiveId: userID,
			Message:   message,
		}

		offlineMessages = append(offlineMessages, offlineMessage)
	}

	// 使用互斥锁保证并发安全
	mutex.Lock()
	// 将离线消息切片保存到数据库表im_offline_messages中
	MYSQL.DataBase.Table("im_offline_messages").Create(&offlineMessages)
	mutex.Unlock()
}
