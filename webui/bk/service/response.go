package service

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Response 通用响应结构体
type Response struct {
	Code      int         `json:"code"`           // 状态码
	Message   string      `json:"message"`        // 响应消息
	Data      interface{} `json:"data,omitempty"` // 响应数据
	Timestamp int64       `json:"timestamp"`      // 时间戳
}

// PageResponse 分页响应结构体
type PageResponse struct {
	Success   bool        `json:"success"`             // 请求是否成功
	Code      int         `json:"code"`                // 状态码
	Message   string      `json:"message"`             // 响应消息
	Data      interface{} `json:"data,omitempty"`      // 响应数据
	Total     int64       `json:"total,omitempty"`     // 总记录数
	Page      int         `json:"page,omitempty"`      // 当前页码
	PageSize  int         `json:"page_size,omitempty"` // 每页大小
	Timestamp int64       `json:"timestamp"`           // 时间戳
}

type Status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// 状态定义
var (
	CodeSuccess = Status{Code: 200, Message: "成功"} // 成功
	CodeErr     = Status{Code: 500, Message: "失败"} // 失败
)

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:      CodeSuccess.Code,
		Message:   CodeSuccess.Message,
		Data:      data,
		Timestamp: time.Now().Unix(),
	})
}

// Error 错误响应
func Error(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, Response{
		Code:      CodeErr.Code,
		Message:   CodeErr.Message,
		Data:      err.Error(),
		Timestamp: time.Now().Unix(),
	})
}
