package auth

import (
	"errors"
	"strings"
	"testing"
)

func TestHashPasswordAndVerifyPassword(t *testing.T) {
	const password = "123456"

	hashed, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() 返回错误：%v", err)
	}
	if !strings.HasPrefix(hashed, "$2") {
		t.Fatalf("HashPassword() 结果不像 bcrypt 哈希：%q", hashed)
	}
	if hashed == password {
		t.Fatal("HashPassword() 不能原样返回明文密码")
	}
	if !VerifyPassword(password, hashed) {
		t.Fatal("VerifyPassword() 对正确密码返回 false")
	}
	if VerifyPassword("wrong-password", hashed) {
		t.Fatal("VerifyPassword() 对错误密码返回 true")
	}
	if VerifyPassword("", hashed) {
		t.Fatal("VerifyPassword() 对空密码返回 true")
	}
	if VerifyPassword(password, "") {
		t.Fatal("VerifyPassword() 对空哈希返回 true")
	}
}

func TestHashPasswordRejectsPasswordLongerThanBcryptLimit(t *testing.T) {
	longPassword := strings.Repeat("a", MaxPasswordBytes+1)
	if _, err := HashPassword(longPassword); !errors.Is(err, ErrPasswordTooLong) {
		t.Fatalf("HashPassword() 对超长密码返回 %v，期望 ErrPasswordTooLong", err)
	}

	limitPassword := strings.Repeat("a", MaxPasswordBytes)
	if _, err := HashPassword(limitPassword); err != nil {
		t.Fatalf("HashPassword() 对 %d 字节密码返回错误：%v", MaxPasswordBytes, err)
	}

	if VerifyPassword(longPassword, "$2a$10$invalid") {
		t.Fatal("VerifyPassword() 对超长密码返回 true")
	}
}
