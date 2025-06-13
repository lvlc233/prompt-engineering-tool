package base


import (

	"github.com/google/uuid"
)

// 生成UUID的辅助函数
func GenerateUUID() string {
    return uuid.New().String()
}