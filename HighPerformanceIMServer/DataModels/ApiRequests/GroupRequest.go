package ApiRequests

type CreateGroupRequest struct {
	UserId     int64    `json:"user_id"`                                // 用户ID
	Name       string   `json:"name" validate:"required,min=2,max=20"`  // 群名称，必填，长度在2到20之间
	Info       string   `json:"info" validate:"required,min=2,max=255"` // 群介绍，必填，长度在2到255之间
	Avatar     string   `json:"avatar" validate:"required"`             // 群头像，必填
	Password   string   `json:"password"`                               // 群密码
	Theme      string   `json:"theme" validate:"required"`              // 群主题，必填
	IsPwd      int      `json:"is_pwd"`                                 // 是否有密码
	SelectUser []string `form:"select_user[]"`                          // 选中的用户列表
}

type CreateUserToGroupRequest struct {
	UserId  []string `json:"select_user[]" validate:"required"` // 用户ID列表，必填
	GroupId int64    `json:"group_id" validate:"required"`      // 群组ID，必填
	Type    int      `json:"type" validate:"required"`          // 类型，必填
}

type InviteUserRequest struct {
	GroupId int64 `json:"group_id" validate:"required"` // 群组ID，必填
	UserId  int64 `json:"user_id" validate:"required"`  // 用户ID，必填
}
