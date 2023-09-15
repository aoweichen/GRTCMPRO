package AUTHGRPC

import (
	"HighPerformanceIMServerAuth/Configs"
	"HighPerformanceIMServerAuth/DAO/MySQL"
	"HighPerformanceIMServerAuth/DataModels/EmailCode"
	"HighPerformanceIMServerAuth/DataModels/IMUsers"
	EmailCodeHtmlTemlate "HighPerformanceIMServerAuth/HtmlTemplate/HTML/EmailCodeHtml"
	"HighPerformanceIMServerAuth/Packages/Date"
	"HighPerformanceIMServerAuth/Packages/EmailServices"
	"HighPerformanceIMServerAuth/Packages/Hash"
	"HighPerformanceIMServerAuth/Packages/JWT"
	"HighPerformanceIMServerAuth/Packages/SetAvatar"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type IMAuthGRPCService struct {
}

func (S *IMAuthGRPCService) mustEmbedUnimplementedIMAuthServiceServer() {}

// Login GRPC服务端登录逻辑
func (S *IMAuthGRPCService) Login(ctx context.Context, request *LoginRequest) (*LoginResponse, error) {
	// 必须的变量
	var user IMUsers.ImUsers
	var response *LoginResponse
	// 查询用户是否存在
	result := MySQL.DataBase.Table("im_users").Where("email=?", request.Email).First(&user)
	if result.RowsAffected == 0 {
		zap.S().Errorf("邮箱 %#v 未注册账号!", request.Email)
		response = &LoginResponse{}
		return response, errors.New("邮箱 " + request.Email + " 未注册账号")
	}
	//	密码错误的情况进行处理
	if !Hash.CheckPassword(request.Password, user.Password) {
		zap.S().Errorln("用户输入的密码错误!")
		response = &LoginResponse{}
		return response, errors.New("用户输入的密码错误")
	}
	// 更新JWT持续时间
	jwtTokenTimeToLive := Configs.ConfigData.JWT.TokenTimeToLive
	expireAtTime := time.Now().Unix() + jwtTokenTimeToLive
	token := JWT.CreateNewJWT().IssueToken(user.ID, user.Uid, user.Name, user.Email, expireAtTime)
	// 登录成功,返回响应体
	response = &LoginResponse{
		ID:              user.ID,
		UID:             user.Uid,
		Name:            user.Name,
		Email:           user.Email,
		Avatar:          user.Avatar,
		Token:           token,
		ExpireTime:      expireAtTime,
		TokenTimeToLive: jwtTokenTimeToLive,
	}
	zap.S().Infoln("用户: " + request.Email + " 登陆成功")
	return response, nil
}

// Register GRPC服务端注册逻辑
func (S *IMAuthGRPCService) Register(ctx context.Context, request *RegisterRequest) (*RegisterResponse, error) {
	var response *RegisterResponse
	// 验证用户是否存在
	if ok, filed := S.IsUserExits(request.Email, request.Name); ok {
		zap.S().Errorf("%s已经存在了", filed)
		response = &RegisterResponse{}
		return response, errors.New(fmt.Sprintf("%s已经存在了", filed))
	}
	// 创建一个EmailService类型的变量emailService
	var emailService EmailServices.EmailService
	// 检查 邮件 验证码
	if !emailService.CheckCode(request.Email, request.EmailCode, int(request.EmailType)) {
		zap.S().Errorln("邮件验证码不正确")
		response = &RegisterResponse{}
		return response, errors.New("邮件验证码不正确")
	}
	// 注册成功
	S.CreateUser(request.Email, request.Password, request.Name) // 创建用户
	response = &RegisterResponse{
		IsRegisterSuccess: true,
	}
	zap.S().Infoln("注册成功！")
	return response, nil
}

// SendEmailCode GRPC服务端发送验证码逻辑
func (S *IMAuthGRPCService) SendEmailCode(ctx context.Context, request *EmailCodeRequest) (*EmailCodeResponse, error) {
	var response *EmailCodeResponse
	// 检查数据库中是否存在指定的邮箱地址
	ok := S.IsTableFiledExits("email", request.Email, "im_users")

	switch request.EmailType {
	// 如果邮箱类型为注册码
	case int64(EmailCode.RegisteredCode):
		// 返回邮箱已经被注册的响应
		if ok {
			zap.S().Infoln("邮箱: " + request.Email + " 已经被注册了")
			response = &EmailCodeResponse{
				IsSendEmailCodeSuccess: false,
			}
			return response, errors.New("邮箱: " + request.Email + " 已经被注册了")
		}
	case int64(EmailCode.ResetPasswordCode):
		if !ok {
			zap.S().Infoln("邮箱: " + request.Email + " 未被注册")
			response = &EmailCodeResponse{
				IsSendEmailCodeSuccess: false,
			}
			return response, errors.New("邮箱: " + request.Email + " 未被注册")
		}
	}
	// 创建一个EmailService类型的变量emailService
	var emailService EmailServices.EmailService
	// 生成一个邮件验证码
	code := S.CreateEmailCode()
	// 生成邮件内容的HTML模板
	html := EmailCodeHtmlTemlate.EmailCodeHTMLTemplate(Configs.ConfigData.Mail.EmailCodeHtmlTemplateFilePath, code)
	// 发送邮件
	emailServiceSendEmailError := emailService.SendEmail(code, int(request.EmailType), request.Email, Configs.ConfigData.Mail.EmailCodeSubject, html)
	if emailServiceSendEmailError != nil {
		zap.S().Errorf("发送失败邮箱:" + request.Email + "错误日志:" + emailServiceSendEmailError.Error())
		// 返回邮件发送失败的响应
		response = &EmailCodeResponse{
			IsSendEmailCodeSuccess: false,
		}
		return response, errors.New("邮件发送失败,请检查是否是可用邮箱")
	}
	// 返回成功响应
	response = &EmailCodeResponse{
		IsSendEmailCodeSuccess: true,
	}
	zap.S().Infoln("邮件发送成功，请注意查收！")
	return response, nil
}

// IMAuthenticateHandler GRPC服务端鉴权逻辑
func (S *IMAuthGRPCService) IMAuthenticateHandler(ctx context.Context, request *AuthRequest) (*AuthResponse, error) {
	var response *AuthResponse

	// 调用 CreateNewJWT 函数创建一个新的 JWT 实例，然后调用 ParseToken 方法解析 token，并返回解析后的 claims 和解析错误
	claims, jWTCreateNewJWTParseTokenError := JWT.CreateNewJWT().ParseToken(request.Token)
	// 如果解析错误不为空，则表示 token 解析失败
	if jWTCreateNewJWTParseTokenError != nil {
		// 使用 zap 记录错误日志
		zap.S().Errorln(jWTCreateNewJWTParseTokenError)
		// 终止后续中间件和请求处理函数的执行
		response = &AuthResponse{}
		return response, jWTCreateNewJWTParseTokenError
	}

	// 鉴权成功
	zap.S().Infoln("鉴权消息成功")
	response = &AuthResponse{
		IsAuthSuccess: true,
		ID:            claims.ID,
		UID:           claims.UID,
		Name:          claims.Name,
	}
	return response, nil
}

// IsTableFiledExits 检查数据库中是否存在指定的邮箱地址
func (S *IMAuthGRPCService) IsTableFiledExits(filed string, value string, table string) bool {
	// 声明一个变量count，用于存储查询结果的数量
	var count int64
	// 在指定的数据表中查询满足条件的记录数量，并将结果存储到count变量中
	MySQL.DataBase.Table(table).Where(fmt.Sprintf("%s=?", filed), value).Count(&count)
	// 如果count大于0，则表示存在满足条件的记录，返回true
	if count > 0 {
		return true
	} else {
		// 如果count等于0，则表示不存在满足条件的记录，返回false
		return false
	}

}

// CreateEmailCode 函数生成一个4位数的验证码字符串
func (S *IMAuthGRPCService) CreateEmailCode() string {
	// 使用rand.NewSource函数创建一个随机数生成器，并使用当前时间的纳秒级Unix时间戳作为种子
	// 使用rand.New函数创建一个新的随机数生成器
	// 使用Int31n方法生成一个0到9999之间的随机整数，并使用fmt.Sprintf函数将其格式化为四位数的字符串
	return fmt.Sprintf("%04v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(10000))
}

// IsUserExits 函数接受两个参数：email 和 name，并返回一个布尔值和一个字符串
func (S *IMAuthGRPCService) IsUserExits(email string, name string) (bool, string) {

	// 声明一个 ImUsers 类型的变量 user
	var user IMUsers.ImUsers

	// 使用模型的 MYSQLDB 对象执行查询操作，查询条件为 email=? or name =?，其中问号会被实际的参数值替换
	// 如果查询结果不为空，则将查询到的第一个用户赋值给 user 变量
	if result := MySQL.DataBase.Table("im_users").Where("email=? or name =?", email, name).First(&user); result.RowsAffected > 0 {

		// 如果查询到的用户邮箱与传入的 email 参数相等，则返回 true 和字符串 "email"
		if user.Email == email {
			return true, "email"
		}

		// 如果查询到的用户姓名与传入的 name 参数相等，则返回 true 和字符串 "name"
		return true, "name"
	}

	// 如果未查询到符合条件的用户，则返回 false 和空字符串
	return false, ""
}

func (S *IMAuthGRPCService) CreateUser(email string, password string, name string) int64 {
	createdAt := Date.NewDate()
	hashedPassword, saltCryptoHashPasswordError := Hash.SaltCryptoHashPassword(password)
	zap.S().Errorf("加密密码出错，:%#v", saltCryptoHashPasswordError)
	users := IMUsers.ImUsers{
		Email:         email,
		Password:      hashedPassword,
		Name:          name,
		CreatedAt:     createdAt,
		UpdatedAt:     createdAt,
		Avatar:        SetAvatar.GetAvatarBase64Png(fmt.Sprintf("https://api.multiavatar.com/%s.svg", name)),
		LastLoginTime: createdAt,
		Uid:           S.GetUuid(),
		UserJson:      "{这名用户还没有任何简介哦}",
		UserType:      1,
	}
	MySQL.DataBase.Table("im_users").Create(&users)
	return users.ID
}

// GetUuid 函数生成一个UUID（Universally Unique Identifier）字符串
func (S *IMAuthGRPCService) GetUuid() string {
	// 使用uuid.NewV4函数生成一个随机的UUID
	u1, _ := uuid.NewUUID()
	// 将UUID转换为字符串并返回
	return fmt.Sprintf("%s", u1)
}
