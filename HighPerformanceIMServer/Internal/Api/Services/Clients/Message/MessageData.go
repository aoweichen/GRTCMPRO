package Message

import "encoding/json"

type Data struct {
}

type MessageDataInterface interface {
	// GetCreateFriendMessageJson 是一个方法，用于将CreateFriendMessage类型的消息转换为JSON格式的[]byte类型
	// 参数message是一个CreateFriendMessage类型的消息
	// 返回一个[]byte类型的JSON格式消息
	GetCreateFriendMessageJson(message CreateFriendMessage) []byte
}

// GetCreateFriendMessageJson 是一个方法，用于将CreateFriendMessage类型的消息转换为JSON格式的[]byte类型
// 参数message是一个CreateFriendMessage类型的消息
// 返回一个[]byte类型的JSON格式消息
func (MD *Data) GetCreateFriendMessageJson(message CreateFriendMessage) []byte {
	// 使用json.Marshal方法将message转换为JSON格式的[]byte类型
	messageJson, _ := json.Marshal(message)
	return messageJson
}
