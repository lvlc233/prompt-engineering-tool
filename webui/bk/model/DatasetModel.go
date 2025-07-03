package model

//请求结构体
//创建
type CreateDatasetRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

//编辑结构体
type EditDatasetRequest struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	Description    string          `json:"description"`
	EditDataTuples []EditDataTuple `json:"edit_data_tuples"`
}

// 编辑数据集具体数据元组结构体
type EditDataTuple struct {
	ID     string `json:"id"`
	Input  string `json:"input"`
	Output string `json:"output"`
}

// UpdateDatasetRequest 更新数据集请求结构体
type UpdateDatasetRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// CreateDatasetDetailRequest 创建数据集详情请求结构体
type CreateDatasetDetailRequest struct {
	DatasetMapID string  `json:"dataset_map_id" binding:"required"`
	Input        *string `json:"input,omitempty"`
	Target       *string `json:"target,omitempty"`
}

// UpdateDatasetDetailRequest 更新数据集详情请求结构体
type UpdateDatasetDetailRequest struct {
	Input  *string `json:"input,omitempty"`
	Target *string `json:"target,omitempty"`
}
