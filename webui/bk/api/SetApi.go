package api

import (
	"webui-server/service"

	"github.com/gin-gonic/gin"
)

func SetSettingRouter(r *gin.Engine) {

	SettingRouter := r.Group("/setting")
	{

		//获取设置参数
		SettingRouter.GET("/get", service.GetSettings)
		//保存
		SettingRouter.POST("/save", service.SaveSettings)
		//重置
		SettingRouter.GET("/reset", service.ResetSettings)
	}
}

// // Setting 设置结构体
// type Setting struct {
// 	ID          int    `json:"id"`
// 	Key         string `json:"key"`
// 	Value       string `json:"value"`
// 	Description string `json:"description"`
// 	Category    string `json:"category"` // api, model, ui, system
// 	Type        string `json:"type"`     // string, number, boolean, password
// 	Required    bool   `json:"required"`
// 	CreatedAt   string `json:"createdAt"`
// 	UpdatedAt   string `json:"updatedAt"`
// }

// // SettingRequest 设置请求
// type SettingRequest struct {
// 	Key         string `json:"key" binding:"required"`
// 	Value       string `json:"value"`
// 	Description string `json:"description"`
// 	Category    string `json:"category"`
// 	Type        string `json:"type"`
// 	Required    bool   `json:"required"`
// }

// // BatchSettingRequest 批量设置请求
// type BatchSettingRequest struct {
// 	Settings []SettingRequest `json:"settings" binding:"required"`
// }

// // SettingMapManager 设置管理器
// type SettingMapManager struct {
// 	settings map[string]*Setting // 使用key作为map的键
// 	nextID   int
// }

// // NewSettingMapManager 创建新的设置管理器
// func NewSettingMapManager() *SettingMapManager {
// 	manager := &SettingMapManager{
// 		settings: make(map[string]*Setting),
// 		nextID:   1,
// 	}

// 	// 初始化默认设置
// 	manager.initDefaultSettings()
// 	return manager
// }

// // initDefaultSettings 初始化默认设置
// func (sm *SettingMapManager) initDefaultSettings() {
// 	now := time.Now().Format("2006-01-02 15:04:05")

// 	defaultSettings := []Setting{
// // 		{
// 			ID:          sm.nextID,
// 			Key:         "api.base_url",
// 			Value:       "https://api.openai.com/v1",
// 			Description: "API服务的基础URL地址",
// 			Category:    "api",
// 			Type:        "string",
// 			Required:    true,
// 			CreatedAt:   now,
// 			UpdatedAt:   now,
// 		},
// 		{
// 			ID:          sm.nextID + 1,
// 			Key:         "api.model",
// 			Value:       "gpt-3.5-turbo",
// 			Description: "要使用的AI模型名称",
// 			Category:    "api",
// 			Type:        "string",
// 			Required:    true,
// 			CreatedAt:   now,
// 			UpdatedAt:   now,
// 		},
// 		{
// 			ID:          sm.nextID + 2,
// 			Key:         "api.api_key",
// 			Value:       "",
// 			Description: "API访问密钥",
// 			Category:    "api",
// 			Type:        "password",
// 			Required:    true,
// 			CreatedAt:   now,
// 			UpdatedAt:   now,
// 		},
// 		{
// 			ID:          sm.nextID + 3,
// 			Key:         "api.timeout",
// 			Value:       "30",
// 			Description: "API请求超时时间（秒）",
// 			Category:    "api",
// 			Type:        "number",
// 			Required:    false,
// 			CreatedAt:   now,
// 			UpdatedAt:   now,
// 		},
// 		{
// 			ID:          sm.nextID + 4,
// 			Key:         "model.temperature",
// 			Value:       "0.7",
// 			Description: "模型温度参数，控制输出的随机性",
// 			Category:    "model",
// 			Type:        "number",
// 			Required:    false,
// 			CreatedAt:   now,
// 			UpdatedAt:   now,
// 		},
// 		{
// 			ID:          sm.nextID + 5,
// 			Key:         "model.max_tokens",
// 			Value:       "2048",
// 			Description: "最大生成token数量",
// 			Category:    "model",
// 			Type:        "number",
// 			Required:    false,
// 			CreatedAt:   now,
// 			UpdatedAt:   now,
// 		},
// 		{
// 			ID:          sm.nextID + 6,
// 			Key:         "ui.theme",
// 			Value:       "light",
// 			Description: "界面主题",
// 			Category:    "ui",
// 			Type:        "string",
// 			Required:    false,
// 			CreatedAt:   now,
// 			UpdatedAt:   now,
// 		},
// 		{
// 			ID:          sm.nextID + 7,
// 			Key:         "ui.language",
// 			Value:       "zh-CN",
// 			Description: "界面语言",
// 			Category:    "ui",
// 			Type:        "string",
// 			Required:    false,
// 			CreatedAt:   now,
// 			UpdatedAt:   now,
// 		},
// 		{
// 			ID:          sm.nextID + 8,
// 			Key:         "system.auto_save",
// 			Value:       "true",
// 			Description: "是否启用自动保存",
// 			Category:    "system",
// 			Type:        "boolean",
// 			Required:    false,
// 			CreatedAt:   now,
// 			UpdatedAt:   now,
// 		},
// 		{
// 			ID:          sm.nextID + 9,
// 			Key:         "system.log_level",
// 			Value:       "info",
// 			Description: "系统日志级别",
// 			Category:    "system",
// 			Type:        "string",
// 			Required:    false,
// 			CreatedAt:   now,
// 			UpdatedAt:   now,
// 		},
// 	}

// 	for _, setting := range defaultSettings {
// 		sm.settings[setting.Key] = &setting
// 	}
// 	sm.nextID += len(defaultSettings)
// }

// // 全局设置管理器实例
// var settingManager = NewSettingMapManager()

// // GetSettings 获取所有设置
// func GetSettings(c *gin.Context) {
// 	// 获取查询参数
// 	category := c.Query("category")
// 	search := c.Query("search")

// 	// 转换为切片并过滤
// 	var settings []*Setting
// 	for _, setting := range settingManager.settings {
// 		// 分类过滤
// 		if category != "" && setting.Category != category {
// 			continue
// 		}

// 		// 搜索过滤
// 		if search != "" {
// 			if !settingContains(setting.Key, search) &&
// 				!settingContains(setting.Description, search) &&
// 				!settingContains(setting.Category, search) {
// 				continue
// 			}
// 		}

// 		settings = append(settings, setting)
// 	}

// 	// 按分类和键名排序
// 	sortSettings(settings)

// 	c.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"data":    settings,
// 		"total":   len(settings),
// 	})
// }

// // GetSetting 获取单个设置
// func GetSetting(c *gin.Context) {
// 	key := c.Param("key")
// 	if key == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": "设置键不能为空",
// 		})
// 		return
// 	}

// 	setting, exists := settingManager.settings[key]
// 	if !exists {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"success": false,
// 			"message": "设置不存在",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"data":    setting,
// 	})
// }

// CreateSetting 创建设置
// func CreateSetting(c *gin.Context) {
// 	var req SettingRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": "请求参数错误: " + err.Error(),
// 		})
// 		return
// 	}

// 	// 检查设置是否已存在
// 	if _, exists := settingManager.settings[req.Key]; exists {
// 		c.JSON(http.StatusConflict, gin.H{
// 			"success": false,
// 			"message": "设置已存在",
// 		})
// 		return
// 	}

// 	// 验证设置类型
// 	validTypes := []string{"string", "number", "boolean", "password"}
// 	if req.Type != "" && !settingSliceContains(validTypes, req.Type) {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": "无效的设置类型",
// 		})
// 		return
// 	}

// 	// 设置默认值
// 	if req.Type == "" {
// 		req.Type = "string"
// 	}
// 	if req.Category == "" {
// 		req.Category = "system"
// 	}

// 	now := time.Now().Format("2006-01-02 15:04:05")
// 	setting := &Setting{
// 		ID:          settingManager.nextID,
// 		Key:         req.Key,
// 		Value:       req.Value,
// 		Description: req.Description,
// 		Category:    req.Category,
// 		Type:        req.Type,
// 		Required:    req.Required,
// 		CreatedAt:   now,
// 		UpdatedAt:   now,
// 	}

// 	settingManager.settings[req.Key] = setting
// 	settingManager.nextID++

// 	c.JSON(http.StatusCreated, gin.H{
// 		"success": true,
// 		"message": "设置创建成功",
// 		"data":    setting,
// 	})
// }

// // UpdateSetting 更新设置
// func UpdateSetting(c *gin.Context) {
// 	key := c.Param("key")
// 	if key == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": "设置键不能为空",
// 		})
// 		return
// 	}

// 	setting, exists := settingManager.settings[key]
// 	if !exists {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"success": false,
// 			"message": "设置不存在",
// 		})
// 		return
// 	}

// 	var req SettingRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": "请求参数错误: " + err.Error(),
// 		})
// 		return
// 	}

// 	// 验证设置类型
// 	validTypes := []string{"string", "number", "boolean", "password"}
// 	if req.Type != "" && !settingSliceContains(validTypes, req.Type) {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": "无效的设置类型",
// 		})
// 		return
// 	}

// 	// 更新设置信息（不允许修改key）
// 	setting.Value = req.Value
// 	if req.Description != "" {
// 		setting.Description = req.Description
// 	}
// 	if req.Category != "" {
// 		setting.Category = req.Category
// 	}
// 	if req.Type != "" {
// 		setting.Type = req.Type
// 	}
// 	setting.Required = req.Required
// 	setting.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

// 	c.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"message": "设置更新成功",
// 		"data":    setting,
// 	})
// }

// // DeleteSetting 删除设置
// func DeleteSetting(c *gin.Context) {
// 	key := c.Param("key")
// 	if key == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": "设置键不能为空",
// 		})
// 		return
// 	}

// 	_, exists := settingManager.settings[key]
// 	if !exists {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"success": false,
// 			"message": "设置不存在",
// 		})
// 		return
// 	}

// 	delete(settingManager.settings, key)

// 	c.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"message": "设置删除成功",
// 	})
// }

// // BatchUpdateSettings 批量更新设置
// func BatchUpdateSettings(c *gin.Context) {
// 	var req BatchSettingRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": "请求参数错误: " + err.Error(),
// 		})
// 		return
// 	}

// 	var updatedSettings []*Setting
// 	var errors []string

// 	for _, settingReq := range req.Settings {
// 		setting, exists := settingManager.settings[settingReq.Key]
// 		if !exists {
// 			errors = append(errors, "设置不存在: "+settingReq.Key)
// 			continue
// 		}

// 		// 更新设置值
// 		setting.Value = settingReq.Value
// 		if settingReq.Description != "" {
// 			setting.Description = settingReq.Description
// 		}
// 		setting.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
// 		updatedSettings = append(updatedSettings, setting)
// 	}

// 	response := gin.H{
// 		"success": true,
// 		"message": "批量更新完成",
// 		"data": gin.H{
// 			"updated": updatedSettings,
// 			"total":   len(updatedSettings),
// 		},
// 	}

// 	if len(errors) > 0 {
// 		response["errors"] = errors
// 	}

// 	c.JSON(http.StatusOK, response)
// }

// // ResetSettings 重置设置为默认值
// func ResetSettings(c *gin.Context) {
// 	category := c.Query("category")

// 	if category != "" {
// 		// 重置指定分类的设置
// 		for key, setting := range settingManager.settings {
// 			if setting.Category == category {
// 				delete(settingManager.settings, key)
// 			}
// 		}
// 	} else {
// 		// 重置所有设置
// 		settingManager.settings = make(map[string]*Setting)
// 		settingManager.nextID = 1
// 	}

// 	// 重新初始化默认设置
// 	settingManager.initDefaultSettings()

// 	c.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"message": "设置重置成功",
// 	})
// }

// // GetSettingCategories 获取设置分类列表
// func GetSettingCategories(c *gin.Context) {
// 	categories := make(map[string]int)
// 	for _, setting := range settingManager.settings {
// 		categories[setting.Category]++
// 	}

// 	var result []gin.H
// 	for category, count := range categories {
// 		result = append(result, gin.H{
// 			"name":  category,
// 			"count": count,
// 		})
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"data":    result,
// 		"total":   len(result),
// 	})
// }

// // sortSettings 排序设置
// func sortSettings(settings []*Setting) {
// 	if len(settings) <= 1 {
// 		return
// 	}

// 	// 简单的冒泡排序，先按分类，再按键名
// 	for i := 0; i < len(settings)-1; i++ {
// 		for j := 0; j < len(settings)-1-i; j++ {
// 			shouldSwap := false
// 			if settings[j].Category > settings[j+1].Category {
// 				shouldSwap = true
// 			} else if settings[j].Category == settings[j+1].Category {
// 				if settings[j].Key > settings[j+1].Key {
// 					shouldSwap = true
// 				}
// 			}

// 			if shouldSwap {
// 				settings[j], settings[j+1] = settings[j+1], settings[j]
// 			}
// 		}
// 	}
// }

// // settingContains 检查字符串是否包含子字符串（忽略大小写）
// func settingContains(s, substr string) bool {
// 	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
// }

// // settingSliceContains 检查切片是否包含指定元素
// func settingSliceContains(slice []string, item string) bool {
// 	for _, s := range slice {
// 		if strings.EqualFold(s, item) {
// 			return true
// 		}
// 	}
// 	return false
