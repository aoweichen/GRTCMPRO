package IMOfflineMessage

// ImGroupOfflineMessages 定义了一个名为ImGroupOfflineMessages的结构体，表示群组离线消息
type ImGroupOfflineMessages struct {
	Id        int    `gorm:"column:id" json:"id"`                 // 离线消息ID
	Message   string `gorm:"column:message" json:"message"`       // 消息体内容
	SendTime  int    `gorm:"column:send_time" json:"send_time"`   // 消息接收时间
	Status    int8   `gorm:"column:status" json:"status"`         // 消息状态，0表示未推送，1表示已推送
	ReceiveId int    `gorm:"column:receive_id" json:"receive_id"` // 接收消息的用户ID
}
