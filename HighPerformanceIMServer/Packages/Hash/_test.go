package Hash

import (
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestSaltCryptoHashPassword(t *testing.T) {
	password := "testPassword"
	hashedPassword, err := SaltCryptoHashPassword(password)
	if err != nil {
		t.Errorf("加密出错： %#v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		t.Errorf("加密后的密码与原密码不匹配： %#v", err)
	}
}

func TestCheckPassword(t *testing.T) {
	password := "testPassword"
	hashedPassword, _ := SaltCryptoHashPassword(password)

	isValid := CheckPassword(password, hashedPassword)
	if !isValid {
		t.Errorf("验证失败，原密码与哈希密码不匹配")
	}
}
