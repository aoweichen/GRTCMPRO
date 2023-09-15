package ApiRequests

// UpdateFriendRequest 结构体用于更新好友请求的请求参数
type UpdateFriendRequest struct {
	ID     string `json:"id" validate:"required"`                 // 字段 ID 表示好友记录的 ID
	Status int    `json:"status" validate:"required,gte=1,lte=2"` // 字段 Status 表示好友记录的状态，取值范围为 1 和 2
}

// CreateFriendRequest 结构体用于创建好友请求的请求参数
type CreateFriendRequest struct {
	ToId        string `json:"to_id" validate:"required"`       // 字段 ToId 表示目标用户的 ID
	Information string `json:"information" validate:"required"` // 字段 Information 表示好友请求的信息
}

// QueryUserRequest 结构体用于查询用户请求的请求参数
type QueryUserRequest struct {
	Email string `json:"email" validate:"omitempty,email"` // 字段 Email 表示用户的电子邮件地址，可以为空字符串，需符合电子邮件格式
}
