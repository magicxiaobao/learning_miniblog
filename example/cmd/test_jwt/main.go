package main

import (
	"fmt"
	"time"

	"example/internal/pkg/auth"
)

func main() {
	// 初始化JWT配置
	auth.InitToken(auth.TokenConfig{
		SigningKey: "test-secret-key",
		ExpireTime: 24, // 24小时
	})

	// 测试生成令牌
	token, err := auth.GenerateToken(1, "testuser", "admin")
	if err != nil {
		fmt.Printf("生成令牌失败: %v\n", err)
		return
	}

	fmt.Printf("生成的JWT令牌: %s\n\n", token)

	// 测试解析令牌
	claims, err := auth.ParseToken(token)
	if err != nil {
		fmt.Printf("解析令牌失败: %v\n", err)
		return
	}

	// 显示令牌内容
	fmt.Println("令牌解析结果:")
	fmt.Printf("用户ID: %d\n", claims.UserID)
	fmt.Printf("用户名: %s\n", claims.Username)
	fmt.Printf("角色: %s\n", claims.Role)
	fmt.Printf("过期时间: %v\n", claims.ExpiresAt.Time)
	fmt.Printf("颁发时间: %v\n", claims.IssuedAt.Time)
	fmt.Printf("生效时间: %v\n\n", claims.NotBefore.Time)

	// 检查令牌是否过期
	if time.Now().After(claims.ExpiresAt.Time) {
		fmt.Println("令牌已过期")
	} else {
		remaining := claims.ExpiresAt.Time.Sub(time.Now())
		fmt.Printf("令牌剩余有效期: %v\n\n", remaining.Round(time.Second))
	}

	// 测试密码加密
	password := "test-password-123"
	fmt.Printf("原始密码: %s\n", password)

	hashedPassword, err := auth.Encrypt(password)
	if err != nil {
		fmt.Printf("密码加密失败: %v\n", err)
		return
	}

	fmt.Printf("加密后的密码: %s\n", hashedPassword)

	// 测试密码验证
	err = auth.Compare(hashedPassword, password)
	if err != nil {
		fmt.Printf("密码验证失败: %v\n", err)
	} else {
		fmt.Println("密码验证成功!")
	}

	// 测试错误密码
	err = auth.Compare(hashedPassword, "wrong-password")
	if err != nil {
		fmt.Printf("错误密码验证结果: %v (预期失败)\n", err)
	} else {
		fmt.Println("错误密码验证成功! (意外结果)")
	}
}
