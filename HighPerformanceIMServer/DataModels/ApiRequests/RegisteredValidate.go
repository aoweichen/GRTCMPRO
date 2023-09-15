package ApiRequests

type RegisteredForm struct {
	Email          string `validate:"required,email"`                        // 邮箱地址，必填且需要符合邮箱格式
	Name           string `validate:"required"`                              // 姓名，必填
	EmailType      int    `validate:"required,gte=1,lte=2"`                  // 邮箱类型，必填且取值范围为1到2之间（1表示个人邮箱，2表示企业邮箱）
	Password       string `json:"password" validate:"required,min=6,max=20"` // 密码，必填且长度在6到20之间
	PasswordRepeat string `validate:"required,eqcsfield=Password"`           // 重复输入的密码，必填且与密码相同
	Code           string `validate:"required,len=4"`                        // 验证码，必填且长度为4
}
