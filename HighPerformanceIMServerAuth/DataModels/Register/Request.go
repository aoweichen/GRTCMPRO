package Register

type Request struct {
	Email          string `json:"email" validate:"required,email"`
	Name           string `json:"name" validate:"required,name"`
	EmailType      int64  `json:"email_type" validate:"required,gte=1,lte=2"`
	Password       string `json:"password" validate:"required,password"`
	PasswordRepeat string `json:"password_repeat" validate:"required,eqcsfield=Password"`
	EmailCode      string `json:"email_code" validate:"required,len=6"`
}
