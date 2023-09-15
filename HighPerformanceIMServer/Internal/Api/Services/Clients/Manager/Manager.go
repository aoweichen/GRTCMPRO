package Manager

import (
	"HighPerformanceIMServer/Configs"
	"HighPerformanceIMServer/DataModels/ApiRequests"
	"HighPerformanceIMServer/Internal/Api/Services/Clients/GRPCMessage"
	"HighPerformanceIMServer/Internal/Api/Services/Clients/Message"
	"HighPerformanceIMServer/Internal/Api/Services/Dispatch"
	NSQQueue "HighPerformanceIMServer/Internal/Api/Services/MessageQueue/NSQ"
	"HighPerformanceIMServer/Packages/Cache/FriendCache"
	"HighPerformanceIMServer/Packages/CoroutinesPoll"
	"HighPerformanceIMServer/Packages/Enums"
	"HighPerformanceIMServer/Packages/MessageOfflineDao"
	"HighPerformanceIMServer/Packages/Utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/valyala/fastjson"
	"sync"
)

// IMClientManager 是一个结构体类型，用于管理 IMClient 客户端的实例
type IMClientManager struct {
	IMClientMap      map[string]*IMClient // 储存客户端实例的映射表，键为字符串类型的客户端标识符，值为指向 IMClient 实例的指针
	BroadcastChannel chan []byte          // 用于广播消息的通道，类型为 []byte 的通道
	PrivateChannel   chan []byte          // 用于私聊消息的通道，类型为 []byte 的通道
	GroupChannel     chan []byte          // 用于群聊消息的通道，类型为 []byte 的通道
	Register         chan *IMClient       // 用于注册新客户端的通道，类型为指向 IMClient 实例的指针的通道
	CloseClients     chan *IMClient       // 用于关闭客户端的通道，类型为指向 IMClient 实例的指针的通道
	MutexKey         sync.RWMutex         // 读写锁，用于保护并发访问客户端实例的映射表
}

// IMClientManagerInterface 是一个接口类型，定义了 IMClientManager 的管理方法合约
type IMClientManagerInterface interface {
	// SetClient 方法用于设置客户端
	SetClient(client *IMClient)
	// DeleteClient 方法用于删除客户端
	DeleteClient(client *IMClient)
	// StartServer 方法用于启动服务器
	StartServer()
	// SendMessageToSpecifiedClient 方法用于向指定的客户端发送消息
	SendMessageToSpecifiedClient(message []byte, toID string) bool
	// LaunchPrivateMessageNSQ 方法用于处理私聊消息
	LaunchPrivateMessageNSQ(messageBytes []byte)
	// LaunchGroupMessage 方法用于处理群聊消息
	LaunchGroupMessage(messageBytes []byte)
	// LaunchBroadcastMessage 方法用于处理广播消息
	LaunchBroadcastMessage(messageBytes []byte)
	// ConsumingOfflineMessage 方法用于处理离线消息
	ConsumingOfflineMessage(client *IMClient)
	// ConsumingGroupOfflineMessages 方法用于处理群聊离线消息
	ConsumingGroupOfflineMessages(client *IMClient)
	// RadioUserOnlineStatus 方法用于广播用户在线状态
	RadioUserOnlineStatus(client *IMClient)

	// GetOnlineNumbers 方法用于获取在线用户数量
	GetOnlineNumbers() int
	// IsUserOnline 判断用户是否在线
	IsUserOnline(string) bool

	// SendPrivateMessage 方法用于发送私聊消息
	SendPrivateMessage(message ApiRequests.PrivateMessageRequest) (bool, string)
	// SendFriendActionMessage 方法用于发送好友动作消息
	SendFriendActionMessage(message Message.CreateFriendMessage)
}

var (
	IMMessageClientManager = IMClientManager{
		IMClientMap:      make(map[string]*IMClient),
		BroadcastChannel: make(chan []byte),
		PrivateChannel:   make(chan []byte),
		GroupChannel:     make(chan []byte),
		Register:         make(chan *IMClient),
		CloseClients:     make(chan *IMClient),
	}
)
var (
	GroupChannelType   = 2
	PrivateChannelType = 1
)

// GetReceiveIdAndUserMsg 定义了一个名为GetReceiveIdAndUserMsg的函数
// 该函数的作用是从一个字节切片中解析出接收者ID和用户消息
// 函数的参数是一个字节切片msg，表示要解析的消息
// 函数返回两个值：接收者ID和用户消息，分别是字符串类型
func GetReceiveIdAndUserMsg(message []byte) (string, string) {
	// 创建一个fastjson.Parser对象p
	var p fastjson.Parser
	// 使用fastjson.Parser对象p解析字节切片msg，并将解析结果赋值给变量v
	messageFastJson, _ := p.Parse(string(message))
	// 从解析结果v中获取字符串类型的接收者ID，使用fastjson.GetString函数
	receiveID := fastjson.GetString(message, "receive_id")
	// 从解析结果v中获取Object类型的消息，使用v.GetObject方法，并将其转换为字符串类型
	userMsg := messageFastJson.GetObject("msg").String()
	// 返回接收者ID和用户消息
	return receiveID, userMsg
}

// SetClient 是 IMClientManager 类型的方法，用于设置客户端
func (Manager *IMClientManager) SetClient(client *IMClient) {
	// 在对 IMClientMap 进行操作前，先获取互斥锁
	Manager.MutexKey.Lock()
	defer Manager.MutexKey.Unlock()
	// 将客户端添加到 IMClientMap 中，以客户端的 ClientID 作为键
	Manager.IMClientMap[client.ClientID] = client
}

// DeleteClient 是 IMClientManager 类型的方法，用于删除客户端
func (Manager *IMClientManager) DeleteClient(client *IMClient) {
	// 在对 ImClientMap 进行操作前，先获取互斥锁
	Manager.MutexKey.Lock()
	defer Manager.MutexKey.Unlock()
	// 关闭客户端的连接
	client.WSClose()
	// 从 ImClientMap 中删除指定客户端
	delete(Manager.IMClientMap, client.ClientID)
}

// StartServer 定义了一个名为StartServer的函数，该函数属于IMClientManager结构体的方法
func (Manager *IMClientManager) StartServer() {
	// 使用无限循环，不断监听select语句中的各个case
	for {
		// 使用select语句，监听多个channel的操作
		select {
		// 如果有新的客户端注册，从IMMessageClientManager的Register channel接收client
		case client := <-IMMessageClientManager.Register:
			// 调用Manager的SetClient方法，将client加入到IMClientMap中
			Manager.SetClient(client)
			// 调用Manager的ConsumingOfflineMessage方法，处理client的离线消息
			Manager.ConsumingOfflineMessage(client)
			// 调用Manager的ConsumingGroupOfflineMessage方法，处理client的群组离线消息
			Manager.ConsumingGroupOfflineMessages(client)
			// 如果有客户端关闭，从IMMessageClientManager的CloseClients channel接收client
		case client := <-IMMessageClientManager.CloseClients:
			// 调用Manager的DeleteClient方法，从IMClientMap中删除client
			Manager.DeleteClient(client)
			// 如果有私聊消息，从IMMessageClientManager的PrivateChannel channel接收message
		case message := <-IMMessageClientManager.PrivateChannel:
			// 使用协程池CoroutinesPoll中的AntPool提交一个函数，该函数调用Manager的LaunchPrivateMessage方法来处理私聊消息
			err := CoroutinesPoll.AntPool.Submit(func() {
				Manager.LaunchPrivateMessageNSQ(message)
			})
			// 调用Utils中的ErrorHandler方法，处理错误
			Utils.ErrorHandler(err)
			// 如果有群聊消息，从IMMessageClientManager的GroupChannel channel接收groupMessages
		case groupMessages := <-IMMessageClientManager.GroupChannel:
			// 使用协程池CoroutinesPoll中的AntPool提交一个函数，该函数调用Manager的LaunchGroupMessage方法来处理群聊消息
			err := CoroutinesPoll.AntPool.Submit(func() {
				Manager.LaunchGroupMessage(groupMessages)
			})
			// 调用Utils中的ErrorHandler方法，处理错误
			Utils.ErrorHandler(err)
			// 如果有广播消息，从IMMessageClientManager的BroadcastChannel channel接收publicMessages
		case publicMessages := <-IMMessageClientManager.BroadcastChannel:
			// 使用协程池CoroutinesPoll中的AntPool提交一个函数，该函数调用Manager的LaunchBroadcastMessage方法来处理广播消息
			err := CoroutinesPoll.AntPool.Submit(func() {
				Manager.LaunchBroadcastMessage(publicMessages)
			})
			// 调用Utils中的ErrorHandler方法，处理错误
			Utils.ErrorHandler(err)
		}
	}
}

// SendMessageToSpecifiedClient 定义了IMClientManager结构体的SendMessageToSpecifiedClient方法
// 该方法用于向指定客户端发送消息
func (Manager *IMClientManager) SendMessageToSpecifiedClient(message []byte, toID string) bool {
	// 判断指定的toID是否存在于IMClientMap中
	if data, ok := Manager.IMClientMap[toID]; ok {
		// 如果存在，将消息内容发送到对应客户端的Send通道
		data.Send <- message
		return true
	} else {
		// 如果不存在，返回false
		return false
	}
}

// LaunchPrivateMessageNSQ 定义了IMClientManager结构体的LaunchPrivateMessageNSQ方法
// 该方法用于处理私人消息的NSQ发送
func (Manager *IMClientManager) LaunchPrivateMessageNSQ(messageBytes []byte) {
	// 调用GetReceiveIdAndUserMsg函数获取接收者ID和用户消息内容
	receiveID, userMessage := GetReceiveIdAndUserMsg(messageBytes)

	// 判断接收者是否在线
	if client, ok := Manager.IMClientMap[receiveID]; ok {
		// 如果接收者在线，将用户消息内容发送到接收者的Send通道
		client.Send <- []byte(userMessage)
	} else {
		// 如果接收者不在线，调用NSQQueue的ProducerQueue的SendMessage方法发送用户消息内容
		NSQQueue.ProducerQueue.SendMessage([]byte(userMessage))
	}
}

// LaunchGroupMessage 定义了IMClientManager结构体的LaunchGroupMessage方法
// 该方法用于发送群组消息
func (Manager *IMClientManager) LaunchGroupMessage(messageBytes []byte) {
	// 调用GetReceiveIdAndUserMsg函数获取接收者ID和用户消息内容
	receiveID, userMessage := GetReceiveIdAndUserMsg(messageBytes)

	// 判断接收者是否在线
	if client, ok := Manager.IMClientMap[receiveID]; ok {
		// 如果接收者在线，将用户消息内容发送到接收者的Send通道
		client.Send <- []byte(userMessage)
	} else {
		// 如果接收者不在线，调用NSQQueue的ProducerQueue的SendGroupMessage方法发送用户消息内容
		NSQQueue.ProducerQueue.SendGroupMessage([]byte(userMessage))
	}
}

// LaunchBroadcastMessage 定义了IMClientManager结构体的LaunchBroadcastMessage方法
// 该方法用于发送广播消息
func (Manager *IMClientManager) LaunchBroadcastMessage(messageBytes []byte) {
	// 创建fastjson.Parser对象用于解析消息内容
	var messageFastJsonParser fastjson.Parser
	// 使用fastjson解析器解析消息内容
	messageFastJson, _ := messageFastJsonParser.Parse(string(messageBytes))

	// 从消息中获取消息码msg_code
	msgCode, _ := messageFastJson.Get("msg_code").Int()
	var receiveID string
	// 根据消息码判断是创建消息还是其他类型的消息，并从消息中获取接收者ID
	if msgCode == Enums.WsCreate {
		receiveID = messageFastJson.Get("to_id").String()
	} else {
		receiveID = messageFastJson.Get("from_id").String()
	}

	// 判断接收者是否在线
	if client, ok := Manager.IMClientMap[receiveID]; ok {
		// 如果接收者在线，将消息内容发送到接收者的Send通道
		client.Send <- messageBytes
	}
}

// ConsumingOfflineMessage 定义了IMClientManager结构体的ConsumingOfflineMessage方法
// 该方法用于消费离线消息
func (Manager *IMClientManager) ConsumingOfflineMessage(client *IMClient) {
	// 从离线消息数据库中获取指定客户端的私人离线消息列表
	messageList := MessageOfflineDao.OfflineMessage.PullPrivateOfflineMessage(client.ClientID)

	// 遍历离线消息列表，将消息逐个发送给客户端
	for _, message := range messageList {
		_ = client.Socket.WriteMessage(websocket.TextMessage, []byte(message.Message))
	}

	// 如果离线消息列表不为空，更新离线消息的状态为已消费
	if len(messageList) > 0 {
		MessageOfflineDao.OfflineMessage.UpdateOfflineMessageStatus(client.ClientID, PrivateChannelType)
	}
}

// ConsumingGroupOfflineMessages 定义了IMClientManager结构体的ConsumingGroupOfflineMessages方法
// 该方法用于消费群组离线消息
func (Manager *IMClientManager) ConsumingGroupOfflineMessages(client *IMClient) {
	// 从离线消息数据库中获取指定客户端的群组离线消息列表
	messageList := MessageOfflineDao.OfflineMessage.PullGroupOfflineMessage(client.ClientID)

	// 遍历离线消息列表，将消息逐个发送给客户端
	for _, message := range messageList {
		_ = client.Socket.WriteMessage(websocket.TextMessage, []byte(message.Message))
	}

	// 如果离线消息列表不为空，更新离线消息的状态为已消费
	if len(messageList) > 0 {
		MessageOfflineDao.OfflineMessage.UpdateOfflineMessageStatus(client.ClientID, GroupChannelType)
	}
}

// RadioUserOnlineStatus 定义了IMClientManager结构体的RadioUserOnlineStatus方法
// 该方法用于广播用户上线状态
func (Manager *IMClientManager) RadioUserOnlineStatus(client *IMClient) {
	// 从好友缓存中获取指定客户端的好友列表
	if imFriends, err := FriendCache.CacheClients.GetFriendCache(client.ClientID); err == nil {
		// 遍历好友列表
		for _, imFriend := range *imFriends {
			// 判断好友是否在线
			if friendClient, ok := Manager.IMClientMap[imFriend.Uid]; ok {
				// 如果好友在线，向好友客户端发送用户上线的消息
				_ = friendClient.Socket.WriteMessage(websocket.TextMessage,
					[]byte(fmt.Sprintf(`{"code":200,"message":"用户上线了"',"fo_id":%d}`, int(imFriend.ToId))))
			}
		}
	}
}

// GetOnlineNumbers 定义了一个名为GetOnlineNumbers的函数，该函数属于IMClientManager结构体的方法
// 该函数的作用是获取在线用户数量
// 函数的参数是一个指向IMClientManager结构体的指针Manager
func (Manager *IMClientManager) GetOnlineNumbers() int {
	// 使用len函数获取Manager的IMClientMap字段的长度，即在线用户的数量
	return len(Manager.IMClientMap)
}

// SendPrivateMessage 定义了IMClientManager结构体的SendPrivateMessage方法
// 该方法用于发送私人消息
func (Manager *IMClientManager) SendPrivateMessage(message ApiRequests.PrivateMessageRequest) (bool, string) {
	// 检查服务是否开启
	if Configs.ConfigData.Server.ClusterOpen {
		// 获取私人聊天消息的字符串表示
		messageString := messageHandler.GetPrivateChatMessages(message, true)

		// 创建DispatchService对象
		var dispatchService Dispatch.Service

		// 判断消息的接收者是否分配到调度节点
		if ok, node := dispatchService.IsDispatchNode(Utils.Int64ToString(message.ToID)); ok && node != "" {
			// 创建GRPCMessageService对象
			var messageClient GRPCMessage.MessageService

			// 使用GRPCMessageService对象发送消息到指定的调度节点
			messageClient.SendGRPCMessage(messageString, node)
			return true, "消息投递成功"
		}
	}

	// 获取私人聊天消息的字符串表示
	messageString := messageHandler.GetPrivateChatMessages(message, true)

	// 根据消息的频道类型进行处理
	switch message.ChannelType {
	case 1:
		// 如果频道类型为1，尝试将消息发送给指定客户端，如果发送失败，则将消息发送到NSQ的消息队列
		if !IMMessageClientManager.SendMessageToSpecifiedClient([]byte(messageString), Utils.Int64ToString(message.ToID)) {
			NSQQueue.ProducerQueue.SendMessage([]byte(messageString))
		}
	case 2:
		// 如果频道类型为2，尝试将消息发送给指定客户端，如果发送失败，则将消息发送到NSQ的群组消息队列
		if !IMMessageClientManager.SendMessageToSpecifiedClient([]byte(messageString), Utils.Int64ToString(message.ToID)) {
			NSQQueue.ProducerQueue.SendGroupMessage([]byte(messageString))
		}
	default:
		// 不处理其他频道类型
	}

	return true, "Success"
}

// SendFriendActionMessage 定义了IMClientManager结构体的SendFriendActionMessage方法
// 该方法用于发送好友动作消息
func (Manager *IMClientManager) SendFriendActionMessage(message Message.CreateFriendMessage) {
	// 将消息结构体转换为JSON格式的字节数组
	messageJson, _ := json.Marshal(message)
	// 调用SendMessageToSpecifiedClient方法将消息JSON发送给指定客户端
	Manager.SendMessageToSpecifiedClient(messageJson, Utils.Int64ToString(message.ToID))
}

// IsUserOnline 是 IMClientManager 类型的方法，用于判断用户是否在线
func (Manager *IMClientManager) IsUserOnline(toID string) bool {
	// 判断 toID 是否在 IMClientMap 中
	_, ok := Manager.IMClientMap[toID]
	// 如果 toID 在 IMClientMap 中，则返回 true；否则返回 false
	return ok
}
