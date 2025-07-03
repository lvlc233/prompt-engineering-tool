package util

import (
	"crypto/rand"
	"math/big"
	"strings"
	"time"

	"github.com/google/uuid"
)

// GenerateRandomString 生成指定长度的随机字符串
// length: 字符串长度
// includeNumbers: 是否包含数字
// includeSymbols: 是否包含特殊符号
func GenerateRandomString(length int, includeNumbers bool, includeSymbols bool) (string, error) {
	if length <= 0 {
		return "", nil
	}

	// 基础字符集（字母）
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// 根据参数添加数字
	if includeNumbers {
		charset += "0123456789"
	}

	// 根据参数添加特殊符号
	if includeSymbols {
		charset += "!@#$%^&*()_+-=[]{}|;:,.<>?"
	}

	var result strings.Builder
	result.Grow(length)

	for i := 0; i < length; i++ {
		// 使用加密安全的随机数生成器
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result.WriteByte(charset[num.Int64()])
	}

	return result.String(), nil
}

// GenerateSimpleRandomString 生成简单的随机字符串（仅包含字母和数字）
func GenerateSimpleRandomString(length int) (string, error) {
	return GenerateRandomString(length, true, false)
}

// GenerateUUID 生成UUID字符串
func GenerateUUID() string {
	return uuid.New().String()
}

// GenerateShortID 生成短ID（基于时间戳和随机数）
func GenerateShortID() (string, error) {
	// 使用时间戳的后6位
	timestamp := time.Now().UnixNano() % 1000000

	// 生成4位随机字符串
	randomPart, err := GenerateSimpleRandomString(4)
	if err != nil {
		return "", err
	}

	return strings.ToUpper(randomPart) + strings.ToUpper(string(rune(timestamp%26+65))) + string(rune(timestamp/26%26+65)), nil
}

// GenerateNumericString 生成纯数字字符串
func GenerateNumericString(length int) (string, error) {
	if length <= 0 {
		return "", nil
	}

	charset := "0123456789"
	var result strings.Builder
	result.Grow(length)

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result.WriteByte(charset[num.Int64()])
	}

	return result.String(), nil
}
