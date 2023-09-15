package IMOfflineMessage

// ImOfflineMessages 定义了一个名为ImOfflineMessages的结构体，表示离线消息
type ImOfflineMessages struct {
	Id        int64  `gorm:"column:id;primaryKey" json:"id"`      // 离线消息ID
	ReceiveId int64  `gorm:"column:receive_id" json:"receive_id"` // 读取消息的用户ID
	Message   string `gorm:"column:message" json:"message"`       // 消息体内容
	SendTime  int    `gorm:"column:send_time" json:"send_time"`   // 消息接收时间
	Status    int    `gorm:"column:status" json:"status"`         // 消息状态，0表示未推送，1表示已推送
}
