package Hash

import (
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// SaltCryptoHashPassword 函数 SaltCryptoHashPassword 用于加密密码并返回加密后的哈希值和错误信息
func SaltCryptoHashPassword(password string) (string, error) {
	// 使用 bcrypt.GenerateFromPassword 函数生成加密后的密码和错误信息
	hashSaltPassword, bcryptGenerateFromPasswordError := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// 如果加密过程中出现错误，则打印错误信息并返回空字符串和错误对象
	if bcryptGenerateFromPasswordError != nil {
		zap.S().Errorf("加密出错： %#v", bcryptGenerateFromPasswordError)
		return "", bcryptGenerateFromPasswordError
	}

	// 如果加密成功，则打印成功信息并返回加密后的哈希值和 nil 错误对象
	zap.S().Info("密码加密成功！")
	return string(hashSaltPassword), nil
}

// CheckPassword 函数 CheckPassword 用于验证输入的密码与存储的哈希密码是否匹配
func CheckPassword(password, hashSaltPassword string) bool {
	// 使用 bcrypt.CompareHashAndPassword 函数比较输入的密码和哈希密码是否匹配
	bcryptCompareHashAndPasswordError := bcrypt.CompareHashAndPassword([]byte(hashSaltPassword), []byte(password))

	// 如果比较过程中出现错误，则打印错误信息并返回 false
	if bcryptCompareHashAndPasswordError != nil {
		zap.S().Errorf("输入的密码错误： %#v", bcryptCompareHashAndPasswordError)
		return false
	} else {
		// 如果比较成功，则打印成功信息并返回 true
		zap.S().Errorln("密码验证成功！")
		return true
	}
}
