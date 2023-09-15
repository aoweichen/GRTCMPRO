package IMGroupMessages

// ImGroupMessages 定义了ImGroupMessages结构体
type ImGroupMessages struct {
	Id              int64   `gorm:"column:id" form:"id"`                               // 消息ID
	Message         string  `gorm:"column:message" form:"message"`                     // 消息实体
	CreatedAt       string  `gorm:"column:created_at" form:"created_at"`               // 添加时间
	Data            string  `gorm:"column:data" form:"data"`                           // 自定义内容
	SendTime        int64   `gorm:"column:send_time" form:"send_time"`                 // 消息添加时间
	MsgType         int     `gorm:"column:msg_type" form:"msg_type"`                   // 消息类型
	MessageId       int64   `gorm:"column:message_id" form:"message_id"`               // 服务端消息ID
	ClientMessageId int64   `gorm:"column:client_message_id" form:"client_message_id"` // 客户端消息ID
	FromId          int64   `gorm:"column:from_id" form:"from_id"`                     // 消息发送者ID
	GroupId         int64   `gorm:"column:group_id" form:"group_id"`                   // 群聊ID
	Users           ImUsers `gorm:"foreignkey:ID;references:FromId"`                   // 关联用户信息
}

// ImUsers 定义了ImUsers结构体
type ImUsers struct {
	ID     int64  `gorm:"column:id;primaryKey" json:"id"` // 用户ID，主键
	Name   string `gorm:"column:name" json:"name"`        // 用户名
	Email  string `gorm:"column:email" json:"email"`      // 邮箱
	Avatar string `gorm:"column:avatar" json:"avatar"`    // 头像
}
