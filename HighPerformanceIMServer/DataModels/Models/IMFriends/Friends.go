package IMFriends

import "HighPerformanceIMServer/DataModels/Models"

type ImFriends struct {
	Models.BaseModel
	Id        int64   `gorm:"column:id" json:"id"`
	FromId    int64   `gorm:"column:from_id" json:"from_id"`
	ToId      int64   `gorm:"column:to_id" json:"to_id"`
	CreatedAt string  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt string  `gorm:"column:created_at" json:"updated_at"`
	Note      string  `gorm:"column:note" json:"note"`
	TopTime   string  `gorm:"column:top_time" json:"top_time"`
	Status    int     `gorm:"column:status" json:"status"` //0.未置顶 1.已置顶
	Uid       string  `gorm:"column:uid" json:"uid"`
	Users     ImUsers `gorm:"foreignKey:ID;references:ToId"`
}

type ImUsers struct {
	ID            int64  `gorm:"column:id;primaryKey" json:"id"`
	Name          string `gorm:"column:name" json:"name"`
	Email         string `gorm:"column:email" json:"email"`
	Avatar        string `gorm:"column:avatar" json:"avatar"`
	Status        int8   `gorm:"column:status" json:"status"`
	Bio           string `gorm:"column:bio" json:"bio"`
	Sex           int8   `gorm:"column:sex" json:"sex"`
	ClientType    int8   `gorm:"column:client_type" json:"client_type"`
	Age           int    `gorm:"column:age" json:"age"`
	LastLoginTime string `gorm:"column:last_login_time" json:"last_login_time"`
	Uid           string `gorm:"column:uid" json:"uid"`
}
