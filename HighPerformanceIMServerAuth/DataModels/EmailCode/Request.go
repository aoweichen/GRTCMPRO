package EmailCode

type EmailValidateCodeRequest struct {
	Email     string `json:"email" validate:"required,email"`
	EmailType int    `json:"email_type" validate:"required,gte=1,tle=2"`
}

var (
	RegisteredCode    = 1
	ResetPasswordCode = 2
)
