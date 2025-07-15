package api

import (
	"webui-server/service"

	"github.com/gin-gonic/gin"
)

var es *service.EvaluationsetService

func SetEvaluationsetManagerRouter(r *gin.Engine) {

	EvaluationsetManagerRouter := r.Group("/evaluationset_manager")
	{
		//获取所有评测集
		EvaluationsetManagerRouter.GET("/", es.GetEvaluationsets)
		//获取评测集详情
		EvaluationsetManagerRouter.GET("/:id", es.GetEvaluationsetDetail)
		//创建评测集
		EvaluationsetManagerRouter.POST("/add", es.CreateEvaluationset)
		//删除评测集
		EvaluationsetManagerRouter.DELETE("/:id", es.DeleteEvaluationset)
	
		//编辑	
		//设置评价标准
		EvaluationsetManagerRouter.POST("/set/criteria",es.SetEvaluationCriteria)
		//设置分数上限
		EvaluationsetManagerRouter.POST("/set/score_cap",es.SetScoreCap)
		
	
		//绑定数据集合
		EvaluationsetManagerRouter.POST("/bind_dataset", es.BindDataset)
		//获取绑定的数据集
		EvaluationsetManagerRouter.GET("/bind_dataset/:id", es.GetBindedDataset)
		//解绑数据集
		EvaluationsetManagerRouter.DELETE("/unbind_dataset/:id", es.UnbindDataset)
		//解绑数据集批次
		EvaluationsetManagerRouter.DELETE("/unbind_dataset_batch", es.UnbindDatasetBatch)
	}
}

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"strconv"
// 	"time"

// 	"github.com/gin-gonic/gin"
// )

// // Evaluation 评测集结构体
// type Evaluation struct {
// 	ID          int    `json:"id"`
// 	Name        string `json:"name"`
// 	Description string `json:"description"`
// 	Type        string `json:"type"`        // 评测类型：accuracy, precision, recall, f1, custom
// 	Metrics     string `json:"metrics"`     // 评测指标配置
// 	CreatedAt   string `json:"createdAt"`
// 	UpdatedAt   string `json:"updatedAt"`
// 	Status      string `json:"status"`      // active, inactive, archived
// }

// // EvaluationRequest 创建/更新评测集请求
// type EvaluationRequest struct {
// 	Name        string `json:"name" binding:"required"`
// 	Description string `json:"description"`
// 	Type        string `json:"type" binding:"required"`
// 	Metrics     string `json:"metrics"`
// 	Status      string `json:"status"`
// }

// // EvaluationMapManager 评测集管理器
// type EvaluationMapManager struct {
// 	evaluations map[int]*Evaluation
// 	nextID      int
// }

// // NewEvaluationMapManager 创建新的评测集管理器
// func NewEvaluationMapManager() *EvaluationMapManager {
// 	return &EvaluationMapManager{
// 		evaluations: make(map[int]*Evaluation),
// 		nextID:      1,
// 	}
// }

// // 全局评测集管理器实例
// var evaluationManager = NewEvaluationMapManager()

// // GetEvaluations 获取所有评测集
// func GetEvaluations(c *gin.Context) {
// 	// 获取查询参数
// 	search := c.Query("search")
// 	sortBy := c.DefaultQuery("sortBy", "time")
// 	sortOrder := c.DefaultQuery("sortOrder", "desc")
// 	status := c.Query("status")
// 	evalType := c.Query("type")

// 	// 转换为切片并过滤
// 	var evaluations []*Evaluation
// 	for _, evaluation := range evaluationManager.evaluations {
// 		// 搜索过滤
// 		if search != "" {
// 			if !evalContains(evaluation.Name, search) && !evalContains(evaluation.Description, search) {
// 				continue
// 			}
// 		}

// 		// 状态过滤
// 		if status != "" && evaluation.Status != status {
// 			continue
// 		}

// 		// 类型过滤
// 		if evalType != "" && evaluation.Type != evalType {
// 			continue
// 		}

// 		evaluations = append(evaluations, evaluation)
// 	}

// 	// 排序
// 	sortEvaluations(evaluations, sortBy, sortOrder)

// 	c.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"data":    evaluations,
// 		"total":   len(evaluations),
// 	})
// }

// // GetEvaluation 获取单个评测集
// func GetEvaluation(c *gin.Context) {
// 	idStr := c.Param("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": "无效的评测集ID",
// 		})
// 		return
// 	}

// 	evaluation, exists := evaluationManager.evaluations[id]
// 	if !exists {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"success": false,
// 			"message": "评测集不存在",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"data":    evaluation,
// 	})
// }

// // CreateEvaluation 创建评测集
// func CreateEvaluation(c *gin.Context) {
// 	var req EvaluationRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": "请求参数错误: " + err.Error(),
// 		})
// 		return
// 	}

// 	// 验证评测类型
// 	validTypes := []string{"accuracy", "precision", "recall", "f1", "custom"}
// 	if !evalSliceContains(validTypes, req.Type) {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": "无效的评测类型",
// 		})
// 		return
// 	}

// 	// 设置默认状态
// 	if req.Status == "" {
// 		req.Status = "active"
// 	}

// 	now := time.Now().Format("2006-01-02 15:04:05")
// 	evaluation := &Evaluation{
// 		ID:          evaluationManager.nextID,
// 		Name:        req.Name,
// 		Description: req.Description,
// 		Type:        req.Type,
// 		Metrics:     req.Metrics,
// 		CreatedAt:   now,
// 		UpdatedAt:   now,
// 		Status:      req.Status,
// 	}

// 	evaluationManager.evaluations[evaluationManager.nextID] = evaluation
// 	evaluationManager.nextID++

// 	c.JSON(http.StatusCreated, gin.H{
// 		"success": true,
// 		"message": "评测集创建成功",
// 		"data":    evaluation,
// 	})
// }

// // UpdateEvaluation 更新评测集
// func UpdateEvaluation(c *gin.Context) {
// 	idStr := c.Param("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": "无效的评测集ID",
// 		})
// 		return
// 	}

// 	evaluation, exists := evaluationManager.evaluations[id]
// 	if !exists {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"success": false,
// 			"message": "评测集不存在",
// 		})
// 		return
// 	}

// 	var req EvaluationRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": "请求参数错误: " + err.Error(),
// 		})
// 		return
// 	}

// 	// 验证评测类型
// 	validTypes := []string{"accuracy", "precision", "recall", "f1", "custom"}
// 	if !evalSliceContains(validTypes, req.Type) {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": "无效的评测类型",
// 		})
// 		return
// 	}

// 	// 更新评测集信息
// 	evaluation.Name = req.Name
// 	evaluation.Description = req.Description
// 	evaluation.Type = req.Type
// 	evaluation.Metrics = req.Metrics
// 	if req.Status != "" {
// 		evaluation.Status = req.Status
// 	}
// 	evaluation.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

// 	c.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"message": "评测集更新成功",
// 		"data":    evaluation,
// 	})
// }

// // DeleteEvaluation 删除评测集
// func DeleteEvaluation(c *gin.Context) {
// 	idStr := c.Param("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": "无效的评测集ID",
// 		})
// 		return
// 	}

// 	_, exists := evaluationManager.evaluations[id]
// 	if !exists {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"success": false,
// 			"message": "评测集不存在",
// 		})
// 		return
// 	}

// 	delete(evaluationManager.evaluations, id)

// 	c.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"message": "评测集删除成功",
// 	})
// }

// // RunEvaluation 运行评测
// func RunEvaluation(c *gin.Context) {
// 	idStr := c.Param("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": "无效的评测集ID",
// 		})
// 		return
// 	}

// 	evaluation, exists := evaluationManager.evaluations[id]
// 	if !exists {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"success": false,
// 			"message": "评测集不存在",
// 		})
// 		return
// 	}

// 	// 模拟评测运行
// 	result := map[string]interface{}{
// 		"evaluationId": id,
// 		"name":         evaluation.Name,
// 		"type":         evaluation.Type,
// 		"status":       "running",
// 		"startTime":    time.Now().Format("2006-01-02 15:04:05"),
// 		"progress":     0,
// 		"results":      nil,
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"message": "评测已开始运行",
// 		"data":    result,
// 	})
// }

// // sortEvaluations 排序评测集
// func sortEvaluations(evaluations []*Evaluation, sortBy, sortOrder string) {
// 	if len(evaluations) <= 1 {
// 		return
// 	}

// 	// 简单的冒泡排序
// 	for i := 0; i < len(evaluations)-1; i++ {
// 		for j := 0; j < len(evaluations)-1-i; j++ {
// 			var shouldSwap bool
// 			switch sortBy {
// 			case "name":
// 				if sortOrder == "asc" {
// 					shouldSwap = evaluations[j].Name > evaluations[j+1].Name
// 				} else {
// 					shouldSwap = evaluations[j].Name < evaluations[j+1].Name
// 				}
// 			case "type":
// 				if sortOrder == "asc" {
// 					shouldSwap = evaluations[j].Type > evaluations[j+1].Type
// 				} else {
// 					shouldSwap = evaluations[j].Type < evaluations[j+1].Type
// 				}
// 			default: // time
// 				if sortOrder == "asc" {
// 					shouldSwap = evaluations[j].CreatedAt > evaluations[j+1].CreatedAt
// 				} else {
// 					shouldSwap = evaluations[j].CreatedAt < evaluations[j+1].CreatedAt
// 				}
// 			}

// 			if shouldSwap {
// 				evaluations[j], evaluations[j+1] = evaluations[j+1], evaluations[j]
// 			}
// 		}
// 	}
// }

// // evalContains 检查字符串是否包含子字符串（忽略大小写）
// func evalContains(s, substr string) bool {
// 	return len(s) >= len(substr) && evalFindSubstring(s, substr)
// }

// // evalFindSubstring 查找子字符串
// func evalFindSubstring(s, substr string) bool {
// 	for i := 0; i <= len(s)-len(substr); i++ {
// 		if s[i:i+len(substr)] == substr {
// 			return true
// 		}
// 	}
// 	return false
// }

// // evalSliceContains 检查切片是否包含指定元素
// func evalSliceContains(slice []string, item string) bool {
// 	for _, s := range slice {
// 		if s == item {
// 			return true
// 		}
// 	}
// 	return false
// }
