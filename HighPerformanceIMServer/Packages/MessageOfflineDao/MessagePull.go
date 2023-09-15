package MessageOfflineDao

import (
	"HighPerformanceIMServer/DataModels/Models/IMOfflineMessage"
	"HighPerformanceIMServer/Internal/DAO/MYSQL"
	"HighPerformanceIMServer/Packages/Date"
	"github.com/golang-module/carbon"
)

// PullPrivateOfflineMessage 定义了一个名为PullPrivateOfflineMessage的方法，属于OfflineMessageDao类型的接收者
// 该方法用于获取指定用户ID的私聊离线消息列表
// 方法接收一个receiveId参数，表示要获取离线消息的用户ID
// 方法返回一个IMOfflineMessage.ImOfflineMessages类型的切片，表示离线消息列表
func (OMD *OfflineMessageDao) PullPrivateOfflineMessage(receiveId string) []IMOfflineMessage.ImOfflineMessages {
	// 创建一个IMOfflineMessage.ImOfflineMessages类型的切片，用于存储离线消息列表
	var list []IMOfflineMessage.ImOfflineMessages

	// 计算15天前的时间戳
	timeStamp := carbon.Parse(Date.NewDate()).SubDays(15).Timestamp()

	// 查询数据库表im_offline_messages，获取符合条件的离线消息列表
	MYSQL.DataBase.Table("im_offline_messages").
		Where("status=0 and receive_id=? and send_time>?", receiveId, timeStamp).
		Find(&list)

	// 返回离线消息列表
	return list
}

// PullGroupOfflineMessage 定义了一个名为PullGroupOfflineMessage的方法，属于OfflineMessageDao类型的接收者
// 该方法用于获取指定用户ID的群聊离线消息列表
// 方法接收一个id参数，表示要获取离线消息的用户ID
// 方法返回一个IMOfflineMessage.ImGroupOfflineMessages类型的切片，表示离线消息列表
func (OMD *OfflineMessageDao) PullGroupOfflineMessage(receiveId string) []IMOfflineMessage.ImGroupOfflineMessages {
	// 创建一个IMOfflineMessage.ImGroupOfflineMessages类型的切片，用于存储离线消息列表
	var list []IMOfflineMessage.ImGroupOfflineMessages

	// 计算15天前的时间戳
	timeStamp := carbon.Parse(Date.NewDate()).SubDays(15).Timestamp()

	// 查询数据库表im_offline_messages，获取符合条件的离线消息列表
	MYSQL.DataBase.Model(&IMOfflineMessage.ImGroupOfflineMessages{}).
		Where("status=0 and receive_id=? and send_time>?", receiveId, timeStamp).
		Find(&list)

	// 返回离线消息列表
	return list
}

// UpdateOfflineMessageStatus UpdatePrivateOfflineMessageStatus 定义了一个名为UpdatePrivateOfflineMessageStatus的方法，属于OfflineMessageDao类型的接收者
// 该方法用于更新指定用户ID的私聊离线消息状态
// 方法接收一个receiveId参数，表示要更新离线消息的用户ID
// 方法接收一个channelType参数，表示消息的渠道类型，1表示私聊，其他表示其他类型
func (OMD *OfflineMessageDao) UpdateOfflineMessageStatus(receiveId string, channelType int) {
	// 根据消息的渠道类型，选择不同的数据模型进行更新操作
	if channelType == 1 {
		// 更新私聊离线消息的状态为已读
		MYSQL.DataBase.Model(&IMOfflineMessage.ImOfflineMessages{}).
			Where("status=0 and receive_id=?", receiveId).
			Updates(map[string]interface{}{"status": 1})
	} else {
		// 更新其他类型离线消息的状态为已读
		MYSQL.DataBase.Model(&IMOfflineMessage.ImGroupOfflineMessages{}).
			Where("status=0 and receive_id=?", receiveId).
			Updates(map[string]interface{}{"status": 1})
	}
}
