package ApiRequests

// LoginForm 定义一个名为 LoginForm 的结构体
type LoginForm struct {
	// 定义一个名为Email的字符串类型的字段，用于存储用户的电子邮件地址
	Email string `json:"email" validate:"required,email"`
	// 定义一个名为Password的字符串类型的字段，用于存储用户的密码
	Password string `json:"password" validate:"required"`
}
