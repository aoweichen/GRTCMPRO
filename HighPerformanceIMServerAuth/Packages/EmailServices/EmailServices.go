package EmailServices

import (
	"HighPerformanceIMServerAuth/Configs"
	"HighPerformanceIMServerAuth/DAO/Redis"
	"HighPerformanceIMServerAuth/DataModels/EmailCode"
	"context"
	"crypto/tls"
	"fmt"
	"go.uber.org/zap"
	"net"
	"net/smtp"
	"time"
)

type EmailServiceInterface interface {
	// SendValidateEmail	发送邮件的方法
	SendValidateEmail(code string, emailType int, email string, subject string, body string) error
	// GetHTMLTemplate  获取html模版内容
	GetHTMLTemplate(text string) []byte
	// GetCacheFix 获取缓存key
	GetCacheFix(email string, emailType int) string
	//	CheckCode 检查验证码
	CheckCode(email string, code string, emailType int) bool
}

type EmailService struct {
}

// SendEmail 是一个发送电子邮件的方法
// 参数：
//   - code: 邮件验证码
//   - emailType: 邮件类型
//   - email: 邮件地址
//   - subject: 邮件主题
//   - body: 邮件正文内容
//
// 返回值：
//   - error: 发送邮件过程中可能发生的错误
func (ES EmailService) SendEmail(code string, emailType int, email string, subject string, body string) error {
	// 创建一个包含邮件头的映射
	htmlHeader := make(map[string]string)
	htmlHeader["From"] = "IM-SERVICE:<" + Configs.ConfigData.Mail.Name + ">"
	htmlHeader["To"] = email
	htmlHeader["Subject"] = subject
	htmlHeader["Content-Type"] = "text/html;chartset=UTF-8"

	// 初始化邮件内容
	message := ""

	// 遍历邮件头映射，将键值对格式化为字符串，并拼接到邮件内容中
	for k, v := range htmlHeader {
		message += fmt.Sprintf("%s:%s\r\n", k, v)
	}

	// 添加邮件正文到邮件内容中
	message += "\r\n" + body

	// 创建 SMTP 认证
	authentication := smtp.PlainAuth("",
		Configs.ConfigData.Mail.Name,
		Configs.ConfigData.Mail.Password,
		Configs.ConfigData.Mail.Host,
	)

	// 使用 TLS 发送邮件
	senMailUsingTLSError := senMailUsingTLS(
		fmt.Sprintf("%s:%d", Configs.ConfigData.Mail.Host, Configs.ConfigData.Mail.Port),
		authentication,
		Configs.ConfigData.Mail.Name,
		[]string{email},
		[]byte(message),
	)

	// 检查邮件发送过程中是否发生错误
	if senMailUsingTLSError != nil {
		zap.S().Errorln("send mail using TLS error:", senMailUsingTLSError)
		return senMailUsingTLSError
	}

	// 将验证码存储到 Redis 缓存中
	Redis.AuthRedisDB.Set(context.Background(), ES.getCacheFix(email, emailType), code, time.Minute*5)

	// 输出日志信息并返回空错误
	zap.S().Info("send email and set key success!")
	return nil
}

// 获取缓存 Key
// getCacheFix 函数定义
func (ES EmailService) getCacheFix(email string, emailType int) string {
	// 根据 emailType 的值来决定返回什么
	switch emailType {
	case EmailCode.RegisteredCode:
		// 如果 emailType 是 REGISTERED_CODE，返回一个由 email 和 REGISTERED_CODE 组成的字符串
		return fmt.Sprintf("%s.%d", email, EmailCode.RegisteredCode)
	case EmailCode.ResetPasswordCode:
		// 如果 emailType 是 RESETPS_CODE，返回一个由 email 和 RESETPS_CODE 组成的字符串
		return fmt.Sprintf("%s.%d", email, EmailCode.ResetPasswordCode)
	default:
		// 如果 emailType 不是这两个值中的任何一个，返回一个由 email 和 REGISTERED_CODE 组成的字符串
		return fmt.Sprintf("%s.%d", email, EmailCode.RegisteredCode)
	}
}

// CheckCode 检测邮箱验证码是否正确
func (ES EmailService) CheckCode(email string, code string, emailType int) bool {
	// 通过调用 getCacheFix 方法获取缓存键值
	cacheFix := ES.getCacheFix(email, emailType)
	//
	// 从 Redis 数据库中获取缓存的值
	redisCmd := Redis.AuthRedisDB.Get(context.Background(), cacheFix)
	//
	// 如果获取到的值为空，则返回 false，表示邮箱不正确
	if val, _ := redisCmd.Result(); val != code {
		return false
	} else {
		// 如果获取到的值与传入的验证码相等，则返回 true，表示邮箱正确
		return true
	}
}

// dialSMTP 返回一个连接的dialSMTP client
// 定义一个函数 dialSMTP，用于建立与 SMTP 服务器的连接
func dialSMTP(addr string) (*smtp.Client, error) {
	// 使用 TLS 协议建立与 SMTP 服务器的连接，并返回错误信息
	conn, dialSMTPError := tls.Dial("tcp", addr, nil)
	if dialSMTPError != nil {
		// 如果连接出错，打印错误信息并返回 nil 和错误对象
		zap.S().Errorln("Dial SMTP Error:", dialSMTPError)
		return nil, dialSMTPError
	}
	// 分解主机端口字符串，获取主机名和端口号
	hosts, _, _ := net.SplitHostPort(addr)
	// 使用已建立的连接和主机名创建一个新的 SMTP 客户端对象，并返回该对象
	return smtp.NewClient(conn, hosts)
}

// 参考net/smtp的func SendMail()
// 使用net.Dial连接tls(ssl)端口时,smtp.NewClient()会卡住且不提示err
// len(to)>1时,to[1]开始提示是密送
// senMailUsingTLS 使用 TLS 发送邮件
func senMailUsingTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	// 创建一个 smtp 客户端
	clientSMTP, dialSMTPError := dialSMTP(addr)
	if dialSMTPError != nil {
		zap.S().Errorln("创建 smtp 客户端错误：")
		return dialSMTPError
	}
	defer func(clientSMTP *smtp.Client) {
		clientSMTPCloseError := clientSMTP.Close()
		if clientSMTPCloseError != nil {
			zap.S().Errorln(clientSMTPCloseError)
		}
	}(clientSMTP)
	// 如果提供了认证信息，则进行身份验证
	if auth != nil {
		if ok, _ := clientSMTP.Extension("AUTH"); ok {
			if clientAuthSMTPError := clientSMTP.Auth(auth); clientAuthSMTPError != nil {
				zap.S().Errorln("AUTH 过程中出错：", clientSMTP)
				return clientAuthSMTPError
			}
		}
	}
	// 设置发件人地址
	if clientMailError := clientSMTP.Mail(from); clientMailError != nil {
		zap.S().Errorln(clientMailError)
		return clientMailError
	}
	// 设置收件人地址列表
	for _, addrs := range to {
		if clientSMTPRcptError := clientSMTP.Rcpt(addrs); clientSMTPRcptError != nil {
			return clientSMTPRcptError
		}
	}
	// 获取数据写入器和关闭写入器的方法
	writeCloser, clientSMTPDataError := clientSMTP.Data()
	if clientSMTPDataError != nil {
		zap.S().Errorln(clientSMTPDataError)
		return clientSMTPDataError
	}
	// 将消息写入数据写入器
	_, writeCloserWriteError := writeCloser.Write(msg)
	if writeCloserWriteError != nil {
		zap.S().Errorln(writeCloserWriteError)
		return writeCloserWriteError
	}
	// 关闭数据写入器
	if writeCloserCloseError := writeCloser.Close(); writeCloserCloseError != nil {
		zap.S().Errorln(writeCloserCloseError)
		return writeCloserCloseError
	}
	zap.S().Infoln("senMailUsingTLS Success!")
	return clientSMTP.Quit()
}
