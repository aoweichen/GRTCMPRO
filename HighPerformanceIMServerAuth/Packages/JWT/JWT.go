package JWT

import (
	"HighPerformanceIMServerAuth/Configs"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"time"
)

type JWT struct {
	SigningKey []byte        // 签名密钥，用于进行JWT的签名和验证
	MaxRefresh time.Duration // 最大刷新时间，表示JWT的有效期
}

// CustomClaims 定义一个CustomClaims结构体，用于存储自定义的声明信息
type CustomClaims struct {
	ID         int64  `json:"id"`          // 用户ID
	UID        string `json:"uid"`         // 用户唯一标识符
	Name       string `json:"name"`        // 用户名
	Email      string `json:"email"`       // 用户邮箱
	ExpireTime int64  `json:"expire_time"` // 过期时间
	jwt.RegisteredClaims
	// StandardClaims 结构体实现了 Claims 接口继承了  Valid() 方法
	// JWT 规定了7个官方字段，提供使用：
	// - iss (issuer)：发布者
	// - sub (subject)：主题
	// - iat (Issued At)：生成签名的时间
	// - exp (expiration time)：签名过期时间
	// - aud (audience)：观众，相当于接受者
	// - nbf (Not Before)：生效时间
	// - jti (JWT ID)：编号
}

// TokenInvalid 定义一个错误类型，表示无法处理此token
var (
	TokenInvalid = errors.New("couldn't handle this token")
)

// CreateNewJWT 是一个函数，它返回一个指向 JWT 结构体的指针
func CreateNewJWT() *JWT {
	// 这里我们创建了一个新的 JWT 结构体，并返回其指针
	return &JWT{
		// 我们使用 ConfigModels.ConfigData.JWT.Secret 作为签名密钥
		SigningKey: []byte(Configs.ConfigData.JWT.Secret),
		// 我们使用 ConfigModels.ConfigData.JWT.TimeToLive 作为最大刷新时间，单位为分钟
		MaxRefresh: time.Duration(Configs.ConfigData.JWT.TokenTimeToLive) * time.Minute,
	}
}

// createToken 创建token
func (J *JWT) createToken(claims CustomClaims) (string, error) {
	// 使用HS256签名方法和传入的claims创建一个JWT对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 对token进行签名，并将签名结果存储在response变量中
	response, tokenSignedStringError := token.SignedString(J.SigningKey)

	// 如果签名过程中出现错误，则打印错误信息并返回response和错误类型
	if tokenSignedStringError != nil {
		zap.S().Errorf("token signed string error: %#v", tokenSignedStringError)
		return response, tokenSignedStringError
	} else {
		// 否则，打印成功信息并返回response和错误类型
		zap.S().Infof("token signed string success!")
		return response, tokenSignedStringError
	}
}

// ParseToken 解析 token
func (J *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	// 调用 jwt.ParseWithClaims 函数解析 JWT 令牌，并返回一个错误（如果有的话）
	token, jwtParseWithClaimsError := jwt.ParseWithClaims(tokenString, &CustomClaims{},
		// 传递一个匿名函数作为参数，该函数返回签名密钥（SigningKey）
		func(token *jwt.Token) (interface{}, error) {
			return J.SigningKey, nil
		})
	// 如果解析过程中出现错误，则打印错误信息并返回错误
	if jwtParseWithClaimsError != nil {
		zap.S().Errorf("jwt ParseWithClaims Error: %#v", jwtParseWithClaimsError)
		return nil, jwtParseWithClaimsError
	}
	// 如果解析成功并且令牌中的声明与预期的自定义声明匹配，则返回声明和nil错误
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	// 如果验证失败，则打印提示信息并返回 nil 和 TokenInvalid 错误
	zap.S().Infoln("JWT token 验证失败")
	return nil, TokenInvalid
}

// RefreshToken 刷新 token
func (J *JWT) RefreshToken(tokenString string) (string, error) {
	// 使用 jwt.ParseWithClaims 函数解析传入的 tokenString，并返回一个 token 对象和一个错误对象
	token, jwtParseWithClaimsError := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 在解析过程中，将 J.SigningKey 作为签名密钥传递给 jwt.ParseWithClaims 函数
		return J.SigningKey, nil
	})

	// 如果解析过程中出现错误，则打印错误信息，并返回空字符串和错误对象
	if jwtParseWithClaimsError != nil {
		zap.S().Errorf("jwt ParseWithClaims Error: %#v", jwtParseWithClaimsError)
		return "", jwtParseWithClaimsError
	}

	// 如果解析成功，并且解析出的 token 中包含自定义声明结构体 CustomClaims 的实例
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		// 将 token 中的过期时间设置为当前时间加上一小时
		claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(1 * time.Hour))
		// 调用 createToken 函数创建一个新的 token，并返回该 token 的字符串表示和 nil 错误对象
		return J.createToken(*claims)
	}

	// 如果解析失败或 token 无效，则打印刷新失败的信息，并返回空字符串和 TokenInvalid 错误对象
	zap.S().Infoln("JWT token 刷新失败")
	return "", TokenInvalid
}

// IssueToken 定义一个名为IssueToken的函数，接收五个参数：id, uid, name, email, expireTime，返回值类型为string
func (J *JWT) IssueToken(id int64, uid string, name string, email string, expireTime int64) string {
	// 构造 claims
	claims := CustomClaims{
		ID:         id,         // 设置id字段
		UID:        uid,        // 设置uid字段
		Name:       name,       // 设置name字段
		Email:      email,      // 设置email字段
		ExpireTime: expireTime, // 设置expireTime字段
		RegisteredClaims: jwt.RegisteredClaims{ // 设置RegisteredClaims字段
			Issuer:    Configs.ConfigData.AuthService.Name,          // 设置发行者为服务器名称
			ExpiresAt: jwt.NewNumericDate(time.Unix(expireTime, 0)), // 设置过期时间为Unix时间戳
			NotBefore: jwt.NewNumericDate(time.Now().Local()),       // 设置生效时间为当前本地时间
			IssuedAt:  jwt.NewNumericDate(time.Now().Local()),       // 设置签发时间为当前本地时间
		},
	}
	// 根据 claims 生成token对象
	token, jWTCreateTokenError := J.createToken(claims)
	if jWTCreateTokenError != nil {
		zap.S().Panicf("根据 claims 生成token对象错误： %#v !", jWTCreateTokenError) // 如果创建token出错，打印错误信息并返回空字符串
		return ""
	}
	zap.S().Infoln("生成 token 成功！") // 打印成功信息
	return token                   // 返回生成的token
}
