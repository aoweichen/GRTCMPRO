package MessageDAO

import (
	"HighPerformanceIMServer/DataModels/ApiRequests"
	"HighPerformanceIMServer/DataModels/Models/IMGroupMessages"
	"HighPerformanceIMServer/DataModels/Models/IMMessages"
	"HighPerformanceIMServer/Internal/DAO/MYSQL"
	"HighPerformanceIMServer/Packages/Date"
	"HighPerformanceIMServer/Packages/Utils"
)

// Dao MessageDao 定义了MessageDao结构体
type Dao struct {
}

// DaoInterface MessageDaoInterface 定义了MessageDaoInterface接口
type DaoInterface interface {
	CreateMessage(params ApiRequests.PrivateMessageRequest)         // 创建单个私聊消息
	CreateMultipleMessage(params ApiRequests.PrivateMessageRequest) // 创建多个私聊消息
	CreateGroupMessage(params ApiRequests.PrivateMessageRequest)    // 创建群聊消息
}

// CreateMessage 实现了CreateMessage方法，用于创建单个私聊消息
func (*Dao) CreateMessage(params ApiRequests.PrivateMessageRequest) {
	// 创建IMMessages实例并赋值
	message := IMMessages.ImMessages{
		Msg:       params.Message,                       // 消息内容
		CreatedAt: params.SendTime,                      // 发送时间
		FromId:    params.FromID,                        // 发送者ID
		ToId:      params.ToID,                          // 接收者ID
		IsRead:    0,                                    // 是否已读，初始为0
		MsgType:   params.MsgType,                       // 消息类型
		Status:    1,                                    // 消息状态，初始为1
		Data:      Utils.InterfaceToString(params.Data), // 自定义数据转换为字符串
	}
	// 保存消息到数据库
	MYSQL.DataBase.Save(&message)
}

// CreateMultipleMessage 实现了CreateMultipleMessage方法，用于创建多个私聊消息
func (*Dao) CreateMultipleMessage(params ApiRequests.PrivateMessageRequest) {
	// 创建IMMessages实例并赋值
	message := IMMessages.ImMessages{
		Msg:       params.Message,                       // 消息内容
		CreatedAt: params.SendTime,                      // 发送时间
		FromId:    params.FromID,                        // 发送者ID
		ToId:      params.ToID,                          // 接收者ID
		IsRead:    0,                                    // 是否已读，初始为0
		MsgType:   params.MsgType,                       // 消息类型
		Status:    1,                                    // 消息状态，初始为1
		Data:      Utils.InterfaceToString(params.Data), // 自定义数据转换为字符串
	}
	// 保存消息到数据库
	MYSQL.DataBase.Save(&message)
}

// CreateGroupMessage 实现了CreateGroupMessage方法，用于创建群聊消息
func (*Dao) CreateGroupMessage(params ApiRequests.PrivateMessageRequest) {
	// 创建IMGroupMessages实例并赋值
	message := IMGroupMessages.ImGroupMessages{
		Message:   params.Message,                       // 消息内容
		CreatedAt: params.SendTime,                      // 发送时间
		Data:      Utils.InterfaceToString(params.Data), // 自定义数据转换为字符串
		SendTime:  Date.TimeUnix(),                      // 当前时间的时间戳
		MsgType:   params.MsgType,                       // 消息类型
		FromId:    params.FromID,                        // 发送者ID
		GroupId:   params.ToID,                          // 群聊ID
	}
	// 保存消息到数据库
	MYSQL.DataBase.Save(&message)
}
