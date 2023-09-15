package IMMessage

import (
	"HighPerformanceIMServer/DataModels/ApiRequests"
	"HighPerformanceIMServer/DataModels/Models/IMGroupMessages"
	"HighPerformanceIMServer/DataModels/Models/IMGroupUsers"
	"HighPerformanceIMServer/DataModels/Models/IMUser"
	"HighPerformanceIMServer/Internal/Api/Services/Clients/Manager"
	"HighPerformanceIMServer/Internal/Api/Services/Clients/Message"
	NSQQueue "HighPerformanceIMServer/Internal/Api/Services/MessageQueue/NSQ"
	"HighPerformanceIMServer/Internal/DAO/MYSQL"
	"HighPerformanceIMServer/Internal/DAO/MessageDAO"
	"HighPerformanceIMServer/Packages/Date"
	"HighPerformanceIMServer/Packages/Enums"
	"HighPerformanceIMServer/Packages/Utils"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

var (
	botData    = map[string]string{} // 存储指令
	lock       sync.RWMutex
	messageDao MessageDAO.Dao
)

// Users 定义了Users结构体
type Users struct {
	UserId string `json:"user_id"` // 用户ID
}

// ImGroups 定义了ImGroups结构体
type ImGroups struct {
	ID        int64  `gorm:"column:id" json:"id"`                 // 群聊ID
	UserId    int64  `gorm:"column:user_id" json:"user_id"`       // 创建者ID
	Name      string `gorm:"column:name" json:"name"`             // 群聊名称
	CreatedAt string `gorm:"column:created_at" json:"created_at"` // 添加时间
	Info      string `gorm:"column:info" json:"info"`             // 群聊描述
	Avatar    string `gorm:"column:avatar" json:"avatar"`         // 群聊头像
	IsPwd     int8   `gorm:"column:is_pwd" json:"is_pwd"`         // 是否加密，0表示否，1表示是
	Hot       int    `gorm:"column:hot" json:"hot"`               // 群聊热度
}

// ImSessionsMessage 定义了ImSessionsMessage结构体
type ImSessionsMessage struct {
	MsgCode  int      `json:"msg_code"` // 消息代码
	Sessions Sessions `json:"sessions"` // 会话内容
}

// Sessions 定义了Sessions结构体
type Sessions struct {
	Id          int64    `gorm:"column:id;primaryKey" json:"id"`          // 会话表ID
	FromId      int64    `gorm:"column:from_id" json:"from_id"`           // 发送者ID
	ToId        int64    `gorm:"column:to_id" json:"to_id"`               // 接收者ID
	GroupId     int64    `gorm:"column:group_id" json:"group_id"`         // 群组ID
	CreatedAt   string   `gorm:"column:created_at" json:"created_at"`     // 创建时间
	TopStatus   int      `gorm:"column:top_status" json:"top_status"`     // 置顶状态，0表示否，1表示是
	TopTime     string   `gorm:"column:top_time" json:"top_time"`         // 置顶时间
	Note        string   `gorm:"column:note" json:"note"`                 // 备注
	ChannelType int      `gorm:"column:channel_type" json:"channel_type"` // 频道类型，0表示单聊，1表示群聊
	Name        string   `gorm:"column:name" json:"name"`                 // 会话名称
	Avatar      string   `gorm:"column:avatar" json:"avatar"`             // 会话头像
	Status      int      `gorm:"column:status" json:"status"`             // 会话状态，0表示正常，1表示禁用
	Groups      ImGroups `gorm:"foreignKey:ID;references:GroupId"`        // 关联的ImGroups结构体
}

type MessageService struct {
}

type MessageServiceInterface interface {
	// IsUserOnline 判断用户是否在线
	IsUserOnline(id string) bool
	// SendFriendActionMessage 发送-好友申请或者拒绝好友请求
	SendFriendActionMessage(message Message.CreateFriendMessage)
	// SendPrivateMessage 发送私聊消息
	SendPrivateMessage(message ApiRequests.PrivateMessageRequest) (bool, string)
	// SendGroupMessage 发送群聊消息
	SendGroupMessage(message ApiRequests.PrivateMessageRequest) bool
	// SendVideoMessage 发送视频请求
	SendVideoMessage(message ApiRequests.VideoMessageRequest) bool
	// SendChatMessage 机器人
	SendChatMessage(message ApiRequests.PrivateMessageRequest) (bool, string)
	//
	SendGroupSessionMessage(userIds []string, groupId int64)
	//
	SendCreateUserGroupMessage(users []IMUser.ImUsers, message ApiRequests.PrivateMessageRequest,
		name interface{}, actionType int, userIds []string)
}

// SliceMock 定义了SliceMock结构体
// 该结构体用于模拟切片的结构
type SliceMock struct {
	addr uintptr // 切片底层数组的起始地址
	len  int     // 切片当前元素个数
	cap  int     // 切片的容量
}

// InSlice 实现了InSlice方法，用于判断一个字符串是否存在于切片中
func InSlice(items []string, item string) bool {
	// 使用for循环遍历切片中的每个元素
	for _, eachItem := range items {
		// 如果当前元素等于目标字符串，返回true
		if eachItem == item {
			return true
		}
	}
	// 如果遍历完切片仍未找到目标字符串，返回false
	return false
}

// IsUserOnline 定义了ImMessageService结构体的IsUserOnline方法
// 该方法用于判断用户是否在线
func (*MessageService) IsUserOnline(id string) bool {
	// 在AppClient的ImManager的ImClientMap中查找指定id的用户
	if _, ok := Manager.IMMessageClientManager.IMClientMap[id]; ok {
		// 如果找到了用户，则表示用户在线，返回true
		return true
	} else {
		// 如果未找到用户，则表示用户不在线，返回false
		return false
	}
}

// SendFriendActionMessage 实现了SendFriendActionMessage方法，用于发送好友操作消息
func (*MessageService) SendFriendActionMessage(message Message.CreateFriendMessage) {
	Manager.IMMessageClientManager.SendFriendActionMessage(message)
}

// SendPrivateMessage 实现了SendPrivateMessage方法，用于发送私聊消息
func (*MessageService) SendPrivateMessage(message ApiRequests.PrivateMessageRequest) (bool, string) {
	return Manager.IMMessageClientManager.SendPrivateMessage(message)
}

// SendChatMessage SendChatMessage函数是ImMessageService结构体的方法，用于发送聊天消息
func (*MessageService) SendChatMessage(message ApiRequests.PrivateMessageRequest) (bool, string) {
	// 将消息的接收者ID设置为发送者ID，实现自己给自己发消息的功能
	message.ToID = message.FromID
	// 将发送者ID设置为1，表示消息发送方为系统
	message.FromID = 1
	// 调用GetMessage函数对消息内容进行处理
	message.Message = GetMessage(message.Message)
	// 调用messageDao的CreateMessage方法将消息保存到数据库中
	messageDao.CreateMessage(message)
	// 调用ClientManager的IMMessageClientManager的SendPrivateMessage方法发送私聊消息
	return Manager.IMMessageClientManager.SendPrivateMessage(message)
}

// SendGroupMessage 实现了SendGroupMessage方法，用于发送群聊消息
func (*MessageService) SendGroupMessage(message ApiRequests.PrivateMessageRequest) bool {
	var users []Users

	// 查询群组中的成员列表
	MYSQL.DataBase.Model(&IMGroupUsers.ImGroupUsers{}).
		Where("group_id=?", message.ToID).
		Select([]string{"user_id"}).Find(&users)

	// 遍历成员列表，发送消息给每个成员
	for _, user := range users {
		// 将用户ID转换为int64类型
		message.UserId = Utils.StringToInt64(user.UserId)

		// 将消息转换为JSON格式
		messageJson, _ := json.Marshal(message)

		// 发送消息给指定的客户端
		if isOK := Manager.IMMessageClientManager.SendMessageToSpecifiedClient(messageJson, user.UserId); isOK {
			// 将消息发送到NSQ队列
			NSQQueue.ProducerQueue.SendGroupMessage(messageJson)
		}
	}

	// 创建群聊消息记录
	groupMessage := IMGroupMessages.ImGroupMessages{
		Message:         message.Message,
		CreatedAt:       Date.NewDate(),
		Data:            message.Data,
		SendTime:        Date.TimeUnix(),
		MsgType:         message.MsgType,
		MessageId:       message.MsgId,
		ClientMessageId: message.MsgClientId,
		FromId:          message.FromID,
	}
	MYSQL.DataBase.Model(&IMGroupMessages.ImGroupMessages{}).Create(&groupMessage)

	return true
}

// SendVideoMessage 实现了SendVideoMessage方法，用于发送视频消息
func (*MessageService) SendVideoMessage(message ApiRequests.VideoMessageRequest) bool {
	// 将消息转换为JSON格式
	messageJson, _ := json.Marshal(message)

	// 将接收者ID转换为字符串类型
	receiverID := Utils.Int64ToString(message.ToID)

	// 发送消息给指定的客户端
	return Manager.IMMessageClientManager.SendMessageToSpecifiedClient(messageJson, receiverID)
}

// SendGroupSessionMessage 实现了SendGroupSessionMessage方法，用于向群组会话发送消息
func (*MessageService) SendGroupSessionMessage(userIds []string, groupId int64) {
	// 创建ImSessionsMessage结构体变量
	var message ImSessionsMessage

	// 设置消息类型为WsSession
	message.MsgCode = Enums.WsSession

	// 使用Model.MYSQLDB查询指定群组的会话信息，并将结果保存到message.Sessions中
	MYSQL.DataBase.Table("im_sessions").Where("group_id=?", groupId).Preload("Groups").Find(&message.Sessions)

	// 遍历用户ID列表，发送消息给每个用户
	for _, id := range userIds {
		// 将用户ID转换为int64类型
		message.Sessions.FromId = Utils.StringToInt64(id)

		// 将消息转换为JSON格式
		msg, _ := json.Marshal(message)

		// 从ClientManager.IMMessageClientManager.IMClientMap中获取用户的连接数据
		data, ok := Manager.IMMessageClientManager.IMClientMap[id]

		// 如果找到了用户的连接数据，则向其发送消息
		if ok {
			data.Send <- msg
		}
	}
}

// SendCreateUserGroupMessage 实现了SendCreateUserGroupMessage方法，用于发送创建用户群组的消息
func (*MessageService) SendCreateUserGroupMessage(users []IMUser.ImUsers, message ApiRequests.PrivateMessageRequest,
	name interface{}, actionType int, userIds []string) {
	var username string

	// 遍历用户列表
	for _, value := range users {
		// 判断当前用户是否在目标用户ID列表中
		if InSlice(userIds, Utils.Int64ToString(value.ID)) {
			username = value.Name

			// 遍历用户列表，向每个用户发送消息
			for _, val := range users {
				message.ToID = val.ID

				// 根据不同的操作类型设置消息内容
				if actionType == 1 {
					if value.ID == val.ID {
						message.Message = fmt.Sprintf("%s邀请您加入了群聊", name)
					} else {
						message.Message = fmt.Sprintf("%s邀请%s加入了群聊", name, username)
					}
				} else {
					message.Message = fmt.Sprintf("%s已经移出群聊", val.Name)
				}

				// 转换消息为JSON格式
				msg, _ := json.Marshal(message)

				// 从ClientManager.IMMessageClientManager.IMClientMap中获取用户的连接数据
				data, ok := Manager.IMMessageClientManager.IMClientMap[Utils.Int64ToString(val.ID)]

				// 如果找到了用户的连接数据，则向其发送消息
				if ok {
					data.Send <- msg
				}
			}
		}
	}
}

// GetMessage 获取指定关键字的消息
func GetMessage(key string) string {

	// 如果关键字包含冒号
	if strings.Contains(key, ":") {
		arr := strings.Split(key, ":")

		// 如果切割后的数组长度为2，表示格式正确
		if len(arr) == 2 {
			lock.Lock()
			botData[arr[0]] = arr[1]
			lock.Unlock()
			return "很不错就是这样~"
		}

		// 如果切割后的数组长度大于2，表示格式不正确
		if len(arr) > 2 {
			return "格式不对呀~"
		}
	}

	// 如果关键字在botData中存在，返回对应的值
	if value, ok := botData[key]; ok {
		return value
	} else {
		// 如果关键字在botData中不存在，返回默认提示语
		return "没明白您的意思-暂时还不知道说啥~~~ 你可以通过 xxx:xxx 指令定义消息😊"
	}
}
