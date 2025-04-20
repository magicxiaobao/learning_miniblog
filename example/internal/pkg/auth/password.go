package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// 用于密码加密的参数
type params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

// 默认参数
var defaultParams = &params{
	memory:      64 * 1024, // 64MB
	iterations:  3,         // 3次迭代
	parallelism: 2,         // 2个并行线程
	saltLength:  16,        // 16字节的盐
	keyLength:   32,        // 32字节的密钥
}

// Encrypt 使用Argon2id算法加密密码
func Encrypt(password string) (string, error) {
	// 生成随机盐
	salt := make([]byte, defaultParams.saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// 使用Argon2id算法加密密码
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		defaultParams.iterations,
		defaultParams.memory,
		defaultParams.parallelism,
		defaultParams.keyLength,
	)

	// 格式化加密后的密码
	// 格式：$argon2id$v=19$m=65536,t=3,p=2$<salt>$<hash>
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		defaultParams.memory,
		defaultParams.iterations,
		defaultParams.parallelism,
		b64Salt,
		b64Hash,
	)

	return encodedHash, nil
}

// Compare 验证密码与散列值是否匹配
func Compare(hashedPassword, password string) error {
	// 解析散列字符串
	parts := strings.Split(hashedPassword, "$")
	if len(parts) != 6 {
		return fmt.Errorf("散列格式无效")
	}

	// 检查算法
	if parts[1] != "argon2id" {
		return fmt.Errorf("不支持的算法: %s", parts[1])
	}

	// 解析参数
	var version int
	p := &params{}
	_, err := fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil {
		return err
	}

	_, err = fmt.Sscanf(
		parts[3],
		"m=%d,t=%d,p=%d",
		&p.memory,
		&p.iterations,
		&p.parallelism,
	)
	if err != nil {
		return err
	}

	// 解码盐和散列
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return err
	}
	p.saltLength = uint32(len(salt))

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return err
	}
	p.keyLength = uint32(len(hash))

	// 计算输入密码的散列
	inputHash := argon2.IDKey(
		[]byte(password),
		salt,
		p.iterations,
		p.memory,
		p.parallelism,
		p.keyLength,
	)

	// 安全比较
	if subtle.ConstantTimeCompare(hash, inputHash) == 1 {
		return nil
	}

	return fmt.Errorf("密码不匹配")
}
