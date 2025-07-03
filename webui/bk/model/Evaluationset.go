package model

// ===== 请求结构体 =====

// CreateEvaluationMapRequest 创建评测集请求
type CreateEvaluationMapRequest struct {
	Name               string   `json:"name" binding:"required"`
	SorceCap           float64 `json:"sorce_cap,omitempty"`
	Description        string  `json:"description,omitempty"`
}
//设置评估标准
type SetEvaluationCriteriaRequest struct {
	EvaluationsetID string `json:"evaluationset_id" binding:"required"`
	Criteria        string `json:"criteria" binding:"required"`
}

//设置分数上限
type SetScoreCapRequest struct {
	EvaluationsetID string `json:"evaluationset_id" binding:"required"`
	ScoreCap float64 `json:"score_cap" binding:"required"`
}

//绑定数据集
type BindDatasetRequest struct {
	EvaluationsetID string `json:"evaluationset_id" binding:"required"`
	DatasetIDs       [] string `json:"dataset_ids" binding:"required"`
}

// ===== 响应结构体 =====

// DatasetWithMapping 包含映射ID和数据集信息的结构体
type DatasetWithMapping struct {
	MappingID   string `json:"mapping_id"`
	DatasetID   string `json:"dataset_id"`
	Name        string `json:"name"`
	DataCount   int    `json:"data_count"`
	Description *string `json:"description,omitempty"`
	CreatedAt   string `json:"created_at"`
}
