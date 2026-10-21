package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// MaxPasswordBytes 是 bcrypt 允许的最大密码长度（字节）。
// 中文按 UTF-8 编码每个字符占 3 字节，因此约 24 个字符就会达到上限。
const MaxPasswordBytes = 72

// ErrPasswordTooLong 表示密码超过 bcrypt 的长度上限。
var ErrPasswordTooLong = fmt.Errorf("密码长度不能超过 %d 字节", MaxPasswordBytes)

// HashPassword 使用 bcrypt 对原始密码做加盐哈希，返回可落库的哈希字符串。
func HashPassword(plainPwd string) (string, error) {
	if len(plainPwd) > MaxPasswordBytes {
		return "", ErrPasswordTooLong
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(plainPwd), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("密码哈希失败：%w", err)
	}
	return string(hashed), nil
}

// VerifyPassword 校验原始密码是否与存储的 bcrypt 哈希匹配。
// 密码超过 bcrypt 长度上限或哈希为空时直接返回 false。
func VerifyPassword(plainPwd, hashedPwd string) bool {
	if plainPwd == "" || hashedPwd == "" || len(plainPwd) > MaxPasswordBytes {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd)) == nil
}
