package IMSessions

// ImSessions 定义会话表结构
type ImSessions struct {
	Id          int64    `gorm:"column:id;primaryKey" json:"id"`          // 会话表
	FromId      int64    `gorm:"column:from_id" json:"from_id"`           // 发送者ID
	ToId        int64    `gorm:"column:to_id" json:"to_id"`               // 接收者ID
	GroupId     int64    `gorm:"column:group_id" json:"group_id"`         // 群组ID
	CreatedAt   string   `gorm:"column:created_at" json:"created_at"`     // 创建时间
	TopStatus   int      `gorm:"column:top_status" json:"top_status"`     // 置顶状态 0:否 1:是
	TopTime     string   `gorm:"column:top_time" json:"top_time"`         // 置顶时间
	Note        string   `gorm:"column:note" json:"note"`                 // 备注
	ChannelType int      `gorm:"column:channel_type" json:"channel_type"` // 渠道类型 0:单聊 1:群聊
	Name        string   `gorm:"column:name" json:"name"`                 // 会话名称
	Avatar      string   `gorm:"column:avatar" json:"avatar"`             // 头像
	Status      int      `gorm:"column:status" json:"status"`             // 会话状态 0:正常 1:禁用
	Users       ImUsers  `gorm:"foreignKey:ID;references:ToId"`           // 用户关联
	Groups      ImGroups `gorm:"foreignKey:ID;references:GroupId"`        // 群聊关联
}

// ImUsers 定义用户表结构
type ImUsers struct {
	ID            int64  `gorm:"column:id;foreignKey" json:"id"`                // 用户ID
	Name          string `gorm:"column:name" json:"name"`                       // 用户名
	Email         string `gorm:"column:email" json:"email"`                     // 邮箱
	Avatar        string `gorm:"column:avatar" json:"avatar"`                   // 头像
	Status        int8   `gorm:"column:status" json:"status"`                   // 用户状态 0:离线 1:在线
	Bio           string `gorm:"column:bio" json:"bio"`                         // 用户简介
	Sex           int8   `gorm:"column:sex" json:"sex"`                         // 性别 0:未知 1:男 2:女
	ClientType    int8   `gorm:"column:client_type" json:"client_type"`         // 客户端类型 1:web 2:pc 3:app
	Age           int    `gorm:"column:age" json:"age"`                         // 年龄
	LastLoginTime string `gorm:"column:last_login_time" json:"last_login_time"` // 最后登录时间
}

// ImGroups 定义群聊表结构
type ImGroups struct {
	ID        int64  `gorm:"column:id" json:"id"`                 // 群聊ID
	UserId    int64  `gorm:"column:user_id" json:"user_id"`       // 创建者ID
	Name      string `gorm:"column:name" json:"name"`             // 群聊名称
	CreatedAt string `gorm:"column:created_at" json:"created_at"` // 创建时间
	Info      string `gorm:"column:info" json:"info"`             // 群聊描述
	Avatar    string `gorm:"column:avatar" json:"avatar"`         // 头像
	IsPwd     int8   `gorm:"column:is_pwd" json:"is_pwd"`         // 是否加密 0:否 1:是
	Hot       int    `gorm:"column:hot" json:"hot"`               // 热度
}

// 定义常量
const (
	SessionStatusOk = 0 // 会话状态：正常
	TopStatus       = 0 // 置顶状态：是
	GROUPTYPE       = 2 // 群聊类型
)
