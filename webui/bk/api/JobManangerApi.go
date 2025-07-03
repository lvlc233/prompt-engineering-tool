package api

import (
	"webui-server/service"

	"github.com/gin-gonic/gin"
)

func SetJobMapManagerRouter(r *gin.Engine) {

	JobMapManagerRouter := r.Group("/job_manager")
	{
		//获取所有任务
		JobMapManagerRouter.GET("/get", service.GetJobs)
		//获取任务详情
		JobMapManagerRouter.GET("/detail/:id", service.GetJobDetail)
		//创建任务
		JobMapManagerRouter.POST("/create", service.CreateJob)
		//删除任务
		JobMapManagerRouter.DELETE("/delete/:id", service.DeleteJob)
		//创建任务版本
		JobMapManagerRouter.POST("/version/:id", service.CreateJobVersion)
		//切换任务版本
		JobMapManagerRouter.PUT("/version/:id/:versionId", service.SwitchJobVersion)
		//获取任务版本列表
		JobMapManagerRouter.GET("/versions/:id", service.GetJobVersions)
		//搜索任务
		JobMapManagerRouter.POST("/search", service.SearchJobs)
		//按名称排序
		JobMapManagerRouter.POST("/sort/name/up", service.SortJobsByNameUp)
		JobMapManagerRouter.POST("/sort/name/down", service.SortJobsByNameDown)
		//按时间排序
		JobMapManagerRouter.POST("/sort/time/up", service.SortJobsByTimeUp)
		JobMapManagerRouter.POST("/sort/time/down", service.SortJobsByTimeDown)
	}
}

// // TaskVersion 任务版本结构体
// type TaskVersion struct {
// 	ID        string   `json:"id"`
// 	Message   string   `json:"message"`
// 	Author    string   `json:"author"`
// 	Timestamp string   `json:"timestamp"`
// 	ParentIds []string `json:"parentIds"`
// 	Branch    string   `json:"branch"`
// 	Status    string   `json:"status"` // committed, draft, current
// }

// // Task 任务结构体
// type Task struct {
// 	ID               int           `json:"id"`
// 	Name             string        `json:"name"`
// 	Description      string        `json:"description"`
// 	CreatedAt        string        `json:"createdAt"`
// 	UpdatedAt        string        `json:"updatedAt"`
// 	CurrentVersionId string        `json:"currentVersionId"`
// 	Versions         []TaskVersion `json:"versions"`
// 	Status           string        `json:"status"` // active, paused, completed, archived
// 	Tags             []string      `json:"tags"`
// }

// // TaskRequest 创建/更新任务请求
// type TaskRequest struct {
// 	Name        string   `json:"name" binding:"required"`
// 	Description string   `json:"description"`
// 	Status      string   `json:"status"`
// 	Tags        []string `json:"tags"`
// }

// // VersionRequest 创建版本请求
// type VersionRequest struct {
// 	Message   string   `json:"message" binding:"required"`
// 	Author    string   `json:"author" binding:"required"`
// 	ParentIds []string `json:"parentIds"`
// 	Branch    string   `json:"branch"`
// }

// // TaskMapManager 任务管理器
// type TaskMapManager struct {
// 	tasks  map[int]*Task
// 	nextID int
// }

// // NewTaskMapManager 创建新的任务管理器
// func NewTaskMapManager() *TaskMapManager {
// 	return &TaskMapManager{
// 		tasks:  make(map[int]*Task),
// 		nextID: 1,
// 	}
// }

// // 全局任务管理器实例
// var taskManager = NewTaskMapManager()

// // GetTasks 获取所有任务
// func GetTasks(c *gin.Context) {
// 	// 获取查询参数
// 	search := c.Query("search")
// 	sortBy := c.DefaultQuery("sortBy", "time")
// 	sortOrder := c.DefaultQuery("sortOrder", "desc")
// 	status := c.Query("status")
// 	tag := c.Query("tag")

// 	// 转换为切片并过滤
// 	var tasks []*Task
// 	for _, task := range taskManager.tasks {
// 		// 搜索过滤
// 		if search != "" {
// 			if !taskContains(strings.ToLower(task.Name), strings.ToLower(search)) &&
// 				!taskContains(strings.ToLower(task.Description), strings.ToLower(search)) &&
// 				!taskSliceContains(task.Tags, search) {
// 				continue
// 			}
// 		}

// 		// 状态过滤
// 		if status != "" && task.Status != status {
// 			continue
// 		}

// 		// 标签过滤
// 		if tag != "" && !sliceContains(task.Tags, tag) {
// 			continue
// 		}

// 		tasks = append(tasks, task)
// 	}

// 	// 排序
// 	sortTasks(tasks, sortBy, sortOrder)

// 	c.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"data":    tasks,
// 		"total":   len(tasks),
// 	})
// }

// // GetTask 获取单个任务
// func GetTask(c *gin.Context) {
// 	idStr := c.Param("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": "无效的任务ID",
// 		})
// 		return
// 	}

// 	task, exists := taskManager.tasks[id]
// 	if !exists {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"success": false,
// 			"message": "任务不存在",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"data":    task,
// 	})
// }

// // CreateTask 创建任务
// func CreateTask(c *gin.Context) {
// 	var req TaskRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": "请求参数错误: " + err.Error(),
// 		})
// 		return
// 	}

// 	// 设置默认状态
// 	if req.Status == "" {
// 		req.Status = "active"
// 	}

// 	now := time.Now().Format("2006-01-02 15:04:05")
// 	initialVersionId := generateVersionId()

// 	// 创建初始版本
// 	initialVersion := TaskVersion{
// 		ID:        initialVersionId,
// 		Message:   "初始项目设置",
// 		Author:    "系统",
// 		Timestamp: now,
// 		ParentIds: []string{},
// 		Branch:    "main",
// 		Status:    "current",
// 	}

// 	task := &Task{
// 		ID:               taskManager.nextID,
// 		Name:             req.Name,
// 		Description:      req.Description,
// 		CreatedAt:        now,
// 		UpdatedAt:        now,
// 		CurrentVersionId: initialVersionId,
// 		Versions:         []TaskVersion{initialVersion},
// 		Status:           req.Status,
// 		Tags:             req.Tags,
// 	}

// 	taskManager.tasks[taskManager.nextID] = task
// 	taskManager.nextID++

// 	c.JSON(http.StatusCreated, gin.H{
// 		"success": true,
// 		"message": "任务创建成功",
// 		"data":    task,
// 	})
// }

// // UpdateTask 更新任务
// func UpdateTask(c *gin.Context) {
// 	idStr := c.Param("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": "无效的任务ID",
// 		})
// 		return
// 	}

// 	task, exists := taskManager.tasks[id]
// 	if !exists {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"success": false,
// 			"message": "任务不存在",
// 		})
// 		return
// 	}

// 	var req TaskRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": "请求参数错误: " + err.Error(),
// 		})
// 		return
// 	}

// 	// 更新任务信息
// 	task.Name = req.Name
// 	task.Description = req.Description
// 	if req.Status != "" {
// 		task.Status = req.Status
// 	}
// 	if req.Tags != nil {
// 		task.Tags = req.Tags
// 	}
// 	task.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

// 	c.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"message": "任务更新成功",
// 		"data":    task,
// 	})
// }

// // DeleteTask 删除任务
// func DeleteTask(c *gin.Context) {
// 	idStr := c.Param("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": "无效的任务ID",
// 		})
// 		return
// 	}

// 	_, exists := taskManager.tasks[id]
// 	if !exists {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"success": false,
// 			"message": "任务不存在",
// 		})
// 		return
// 	}

// 	delete(taskManager.tasks, id)

// 	c.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"message": "任务删除成功",
// 	})
// }

// // CreateTaskVersion 创建任务版本
// func CreateTaskVersion(c *gin.Context) {
// 	idStr := c.Param("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": "无效的任务ID",
// 		})
// 		return
// 	}

// 	task, exists := taskManager.tasks[id]
// 	if !exists {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"success": false,
// 			"message": "任务不存在",
// 		})
// 		return
// 	}

// 	var req VersionRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": "请求参数错误: " + err.Error(),
// 		})
// 		return
// 	}

// 	// 设置默认分支
// 	if req.Branch == "" {
// 		req.Branch = "main"
// 	}

// 	// 生成新版本ID
// 	versionId := generateVersionId()

// 	// 创建新版本
// 	newVersion := TaskVersion{
// 		ID:        versionId,
// 		Message:   req.Message,
// 		Author:    req.Author,
// 		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
// 		ParentIds: req.ParentIds,
// 		Branch:    req.Branch,
// 		Status:    "draft",
// 	}

// 	// 添加到任务版本列表
// 	task.Versions = append(task.Versions, newVersion)
// 	task.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

// 	c.JSON(http.StatusCreated, gin.H{
// 		"success": true,
// 		"message": "版本创建成功",
// 		"data":    newVersion,
// 	})
// }

// // SwitchTaskVersion 切换任务版本
// func SwitchTaskVersion(c *gin.Context) {
// 	idStr := c.Param("id")
// 	versionId := c.Param("versionId")

// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": "无效的任务ID",
// 		})
// 		return
// 	}

// 	task, exists := taskManager.tasks[id]
// 	if !exists {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"success": false,
// 			"message": "任务不存在",
// 		})
// 		return
// 	}

// 	// 检查版本是否存在
// 	versionExists := false
// 	for i := range task.Versions {
// 		if task.Versions[i].ID == versionId {
// 			versionExists = true
// 			// 更新版本状态
// 			if task.CurrentVersionId != "" {
// 				// 将当前版本状态改为committed
// 				for j := range task.Versions {
// 					if task.Versions[j].ID == task.CurrentVersionId {
// 						task.Versions[j].Status = "committed"
// 						break
// 					}
// 				}
// 			}
// 			// 设置新的当前版本
// 			task.Versions[i].Status = "current"
// 			break
// 		}
// 	}

// 	if !versionExists {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"success": false,
// 			"message": "版本不存在",
// 		})
// 		return
// 	}

// 	task.CurrentVersionId = versionId
// 	task.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

// 	c.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"message": "版本切换成功",
// 		"data":    task,
// 	})
// }

// // GetTaskVersions 获取任务的所有版本
// func GetTaskVersions(c *gin.Context) {
// 	idStr := c.Param("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": "无效的任务ID",
// 		})
// 		return
// 	}

// 	task, exists := taskManager.tasks[id]
// 	if !exists {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"success": false,
// 			"message": "任务不存在",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"data":    task.Versions,
// 		"total":   len(task.Versions),
// 	})
// }

// // sortTasks 排序任务
// func sortTasks(tasks []*Task, sortBy, sortOrder string) {
// 	if len(tasks) <= 1 {
// 		return
// 	}

// 	// 简单的冒泡排序
// 	for i := 0; i < len(tasks)-1; i++ {
// 		for j := 0; j < len(tasks)-1-i; j++ {
// 			var shouldSwap bool
// 			switch sortBy {
// 			case "name":
// 				if sortOrder == "asc" {
// 					shouldSwap = tasks[j].Name > tasks[j+1].Name
// 				} else {
// 					shouldSwap = tasks[j].Name < tasks[j+1].Name
// 				}
// 			case "status":
// 				if sortOrder == "asc" {
// 					shouldSwap = tasks[j].Status > tasks[j+1].Status
// 				} else {
// 					shouldSwap = tasks[j].Status < tasks[j+1].Status
// 				}
// 			default: // time
// 				if sortOrder == "asc" {
// 					shouldSwap = tasks[j].CreatedAt > tasks[j+1].CreatedAt
// 				} else {
// 					shouldSwap = tasks[j].CreatedAt < tasks[j+1].CreatedAt
// 				}
// 			}

// 			if shouldSwap {
// 				tasks[j], tasks[j+1] = tasks[j+1], tasks[j]
// 			}
// 		}
// 	}
// }

// // generateVersionId 生成版本ID
// func generateVersionId() string {
// 	// 简单的版本ID生成，实际项目中可以使用更复杂的算法
// 	now := time.Now()
// 	return strings.ToLower(now.Format("20060102150405")[:7])
// }

// // taskContains 检查字符串是否包含子字符串（忽略大小写）
// func taskContains(s, substr string) bool {
// 	return len(s) >= len(substr) && taskFindSubstring(s, substr)
// }

// // taskFindSubstring 查找子字符串
// func taskFindSubstring(s, substr string) bool {
// 	for i := 0; i <= len(s)-len(substr); i++ {
// 		if s[i:i+len(substr)] == substr {
// 			return true
// 		}
// 	}
// 	return false
// }

// // taskSliceContains 检查切片是否包含指定元素
// func taskSliceContains(slice []string, item string) bool {
// 	for _, s := range slice {
// 		if s == item {
// 			return true
// 		}
// 	}
// 	return false
// }
