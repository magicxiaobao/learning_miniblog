package auth

import (
	"testing"
)

func TestPasswordEncryptionAndValidation(t *testing.T) {
	// 测试数据
	password := "test-password-123"

	// 测试加密
	hashedPassword, err := Encrypt(password)
	if err != nil {
		t.Fatalf("密码加密失败: %v", err)
	}

	if hashedPassword == "" {
		t.Fatal("加密后的密码为空")
	}

	if hashedPassword == password {
		t.Fatal("加密后的密码与原始密码相同，加密失败")
	}

	// 测试格式
	if len(hashedPassword) < 20 {
		t.Fatalf("加密后的密码太短: %s", hashedPassword)
	}

	// 测试验证正确的密码
	err = Compare(hashedPassword, password)
	if err != nil {
		t.Fatalf("验证正确密码失败: %v", err)
	}

	// 测试验证错误的密码
	err = Compare(hashedPassword, "wrong-password")
	if err == nil {
		t.Fatal("验证错误密码应该失败，但通过了")
	}

	// 测试不同的密码产生不同的哈希
	hashedPassword2, _ := Encrypt(password)
	if hashedPassword == hashedPassword2 {
		t.Fatal("两次加密产生了相同的哈希，安全性存在问题")
	}
}
