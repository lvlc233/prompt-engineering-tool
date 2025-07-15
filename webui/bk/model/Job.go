package model


// ===== 请求结构体 =====

// CreateJobRequest 创建任务请求
type CreateJobRequest struct {
	Name string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

// CreateJobVersionRequest 创建任务版本请求
type CreateJobVersionRequest struct {
	JobID           string `json:"job_id" binding:"required"`
	FatherVersion   string `json:"father_version"`
	Description     string `json:"description"`
	InputPrompt     string `json:"input_prompt" binding:"required"`
	OutputPrompt    string `json:"output_prompt"`
	OptimizeOrientation string `json:"optimize_orientation"`
}

// UpdateJobRequest 更新任务请求
type UpdateJobRequest struct {
	Name string `json:"name" binding:"required"`
}


// JobDetailResponse 任务详情响应
type JobDetailResponse struct {
	JobID               string `json:"job_id"`
	Name                string `json:"name"`
	CreatedAt           string `json:"created_at"`
	Version             string `json:"version"`
	FatherVersion       string `json:"father_version"`
	Description         string `json:"description"`
	InputPrompt         string `json:"input_prompt"`
	OutputPrompt        string `json:"output_prompt"`
	OptimizeOrientation string `json:"optimize_orientation"`
	OptimizedPrompt     string `json:"optimized_prompt"`
	IsExecute           bool   `json:"is_execute"`
	ExecuteDate         string `json:"execute_date"`
}