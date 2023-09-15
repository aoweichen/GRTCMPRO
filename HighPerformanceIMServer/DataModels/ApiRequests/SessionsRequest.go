package ApiRequests

// SessionStore 定义名为SessionStore的结构体
// 用于表示会话的存储请求数据结构
type SessionStore struct {
	Id   int64 `json:"id" validate:"required"`               // 表示会话的ID，必需字段
	Type int   `json:"type" validate:"required,gte=1,lte=2"` // 表示会话的类型，必需字段，取值范围为1到2之间
}

// SessionUpdate 定义名为SessionUpdate的结构体
// 用于表示会话的更新请求数据结构
type SessionUpdate struct {
	TopStatus int    `json:"top_status" validate:"required,gte=0,lte=1"` // 表示会话的置顶状态，必需字段，取值范围为0到1之间
	Note      string `json:"type"`                                       // 表示会话的备注，可选字段
}
