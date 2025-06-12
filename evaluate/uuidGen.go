package evaluate

import (

	"github.com/google/uuid"
)

// 生成UUID的辅助函数
func generateUUID() string {
    return uuid.New().String()
}