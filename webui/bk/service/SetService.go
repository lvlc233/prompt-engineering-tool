package service

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// EnvConfig 环境配置结构体
type EnvConfig struct {
	BaseURL   string `json:"base_url"`
	ModelName string `json:"model_name"`
	APIKey    string `json:"api_key"`
}

// SaveRequest 保存请求结构体
// omitempty 是JSON标签的一个选项，它的作用是：
// 当字段值为零值时，在JSON序列化过程中忽略该字段
type SaveRequest struct {
	BaseURL   *string `json:"base_url,omitempty"`
	ModelName *string `json:"model_name,omitempty"`
	APIKey    *string `json:"api_key,omitempty"`
}

// getEnvFilePath 获取.env文件路径
func getEnvFilePath() string {
	return "..\\..\\.env"
}

// readEnvFile 读取.env文件
func readEnvFile() (*EnvConfig, error) {
	envPath := getEnvFilePath()
	file, err := os.Open(envPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &EnvConfig{}
	scanner := bufio.NewScanner(file)

	// 逐行读取配置
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "base_url":
			config.BaseURL = value
		case "model_name":
			config.ModelName = value
		case "api_key":
			config.APIKey = value
		}
	}

	return config, scanner.Err()
}

// writeEnvFile 写入.env文件
func writeEnvFile(config *EnvConfig) error {
	envPath := getEnvFilePath()

	// 确保目录存在
	dir := filepath.Dir(envPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(envPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "model_name=%s\napi_key=%s\nbase_url=%s\n",
		config.ModelName, config.APIKey, config.BaseURL)
	return err
}

// GetSettings 获取设置参数
func GetSettings(c *gin.Context) {
	config, err := readEnvFile()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "读取配置文件失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    config,
	})
}

// SaveSettings 保存设置
func SaveSettings(c *gin.Context) {
	var req SaveRequest
	//使用c.ShouldBindJSON(&req);进行参数映射
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 读取当前配置
	config, err := readEnvFile()
	if err != nil {
		// 如果文件不存在，创建默认配置
		config = &EnvConfig{
			BaseURL:   "https://api.deepseek.com",
			ModelName: "deepseek-chat",
			APIKey:    "",
		}
	}

	// 更新配置（只更新提供的字段）
	if req.BaseURL != nil {
		config.BaseURL = *req.BaseURL
	}
	if req.ModelName != nil {
		config.ModelName = *req.ModelName
	}
	if req.APIKey != nil {
		config.APIKey = *req.APIKey
	}

	// 保存配置
	if err := writeEnvFile(config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "保存配置文件失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "配置保存成功",
		"data":    config,
	})
}

// ResetSettings 重置设置
func ResetSettings(c *gin.Context) {
	// 默认配置
	defaultConfig := &EnvConfig{
		BaseURL:   "https://api.deepseek.com",
		ModelName: "deepseek-chat",
		APIKey:    "",
	}

	// 保存默认配置
	if err := writeEnvFile(defaultConfig); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "重置配置失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "配置重置成功",
		"data":    defaultConfig,
	})
}
