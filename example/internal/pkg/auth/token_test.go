package auth

import (
	"testing"
	"time"
)

func TestJWTToken(t *testing.T) {
	// 初始化配置
	InitToken(TokenConfig{
		SigningKey: "test-key",
		ExpireTime: 1, // 1小时
	})

	// 测试数据
	userID := uint64(1)
	username := "testuser"
	role := "admin"

	// 测试生成令牌
	token, err := GenerateToken(userID, username, role)
	if err != nil {
		t.Fatalf("无法生成令牌: %v", err)
	}

	if token == "" {
		t.Fatal("生成的令牌为空")
	}

	// 测试解析令牌
	claims, err := ParseToken(token)
	if err != nil {
		t.Fatalf("无法解析令牌: %v", err)
	}

	// 验证令牌内容
	if claims.UserID != userID {
		t.Errorf("用户ID不匹配，预期 %d，得到 %d", userID, claims.UserID)
	}

	if claims.Username != username {
		t.Errorf("用户名不匹配，预期 %s，得到 %s", username, claims.Username)
	}

	if claims.Role != role {
		t.Errorf("角色不匹配，预期 %s，得到 %s", role, claims.Role)
	}

	// 验证过期时间
	now := time.Now()
	if !claims.ExpiresAt.Time.After(now) {
		t.Error("令牌已过期")
	}

	// 测试过期令牌
	// 创建一个过期的令牌配置
	InitToken(TokenConfig{
		SigningKey: "test-key",
		ExpireTime: -1, // -1小时（已过期）
	})

	expiredToken, _ := GenerateToken(userID, username, role)
	_, err = ParseToken(expiredToken)
	if err != ErrTokenExpired {
		t.Errorf("应该返回过期错误，但得到: %v", err)
	}

	// 测试无效令牌
	_, err = ParseToken("invalid.token.string")
	if err != ErrTokenInvalid {
		t.Errorf("应该返回无效令牌错误，但得到: %v", err)
	}
}
