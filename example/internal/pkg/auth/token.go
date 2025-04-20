package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// 定义Token相关的错误
var (
	ErrMissingHeader = errors.New("认证头不存在")
	ErrTokenInvalid  = errors.New("无效的Token")
	ErrTokenExpired  = errors.New("Token已过期")
)

// TokenConfig 包含JWT配置选项
type TokenConfig struct {
	// 用于签名的密钥
	SigningKey string
	// 有效期（以小时为单位）
	ExpireTime int
}

// 默认配置
var defaultConfig = TokenConfig{
	SigningKey: "JnCEIL6CRbnLj6U5PO4nKXJtUZRvDG87",
	ExpireTime: 24, // 24小时
}

// CustomClaims 自定义JWT声明
type CustomClaims struct {
	// 用户信息
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	// 标准JWT声明
	jwt.RegisteredClaims
}

// Init 初始化配置
func InitToken(config TokenConfig) {
	if config.SigningKey != "" {
		defaultConfig.SigningKey = config.SigningKey
	}
	if config.ExpireTime > 0 {
		defaultConfig.ExpireTime = config.ExpireTime
	}
}

// GenerateToken 生成JWT令牌
func GenerateToken(userID uint64, username, role string) (string, error) {
	now := time.Now()
	claims := CustomClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(defaultConfig.ExpireTime) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   username,
		},
	}

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名令牌
	return token.SignedString([]byte(defaultConfig.SigningKey))
}

// ParseToken 解析JWT令牌
func ParseToken(tokenString string) (*CustomClaims, error) {
	// 创建一个CustomClaims实例
	claims := &CustomClaims{}

	// 解析令牌
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("意外的签名方法: %v", token.Header["alg"])
		}
		// 返回密钥
		return []byte(defaultConfig.SigningKey), nil
	})

	// 处理解析错误
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}

	// 验证令牌有效性
	if !token.Valid {
		return nil, ErrTokenInvalid
	}

	return claims, nil
}

// ParseRequest 从HTTP请求中提取并解析JWT令牌
func ParseRequest(c *gin.Context) (*CustomClaims, error) {
	// 从请求头获取令牌
	header := c.Request.Header.Get("Authorization")
	if header == "" {
		return nil, ErrMissingHeader
	}

	// 从Bearer前缀中提取令牌
	var tokenString string
	fmt.Sscanf(header, "Bearer %s", &tokenString)
	if tokenString == "" {
		return nil, ErrTokenInvalid
	}

	// 解析令牌
	return ParseToken(tokenString)
}
