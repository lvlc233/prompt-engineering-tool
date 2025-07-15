package service

import (
	"errors"
	"webui-server/model"
	"webui-server/sql"

	"github.com/gin-gonic/gin"
)

type EvaluationsetService struct{}

var ec = &sql.EvaluationsetCRUD{}
var edmc = &sql.EvaluationsetDatasetMappingCRUD{}

// 获取所有评测集
func (es *EvaluationsetService) GetEvaluationsets(c *gin.Context) {
	evaluationsets, err := ec.GetAllEvaluationsets()
	if err != nil {
		Error(c, err)
		return
	}
	Success(c, evaluationsets)
}

// 获取评测集详情根据评测集ID
func (es *EvaluationsetService) GetEvaluationsetDetail(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		Error(c, errors.New("缺少评测集ID参数"))
		return
	}
	evaluationset, err := ec.GetEvaluationsetByID(id)
	if err != nil {
		Error(c, err)
		return
	}
	if evaluationset == nil {
		Error(c, errors.New("评测集不存在"))
		return
	}
	Success(c, evaluationset)
}

// 创建新的评估集
func (es *EvaluationsetService) CreateEvaluationset(c *gin.Context) {
	var req model.CreateEvaluationMapRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, err)
		return
	}

	err := ec.AddEvaluationset(req.Name, req.Description, req.SorceCap)
	if err != nil {
		Error(c, err)
		return
	}

	Success(c, "评测集创建成功")
}

// 设置评估集的评价标准
func (es *EvaluationsetService) SetEvaluationCriteria(c *gin.Context) {
	var req model.SetEvaluationCriteriaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, err)
		return
	}
	err := ec.SetEvaluationCriteria(req.EvaluationsetID, req.Criteria)
	if err != nil {
		Error(c, err)
		return
	}
	Success(c, "评测集评价标准设置成功")
}

// 设置评估集的分数上限
func (es *EvaluationsetService) SetScoreCap(c *gin.Context) {
	var req model.SetScoreCapRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, err)
		return
	}
	err := ec.SetScoreCap(req.EvaluationsetID, req.ScoreCap)
	if err != nil {
		Error(c, err)
		return
	}
	Success(c, "评测集分数上限设置成功")
}

// 绑定数据集
func (es *EvaluationsetService) BindDataset(c *gin.Context) {
	var req model.BindDatasetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, err)
		return
	}
	err := edmc.AddDatasetMappingByBatch(req.EvaluationsetID, req.DatasetIDs)
	if err != nil {
		Error(c, err)
		return
	}
	Success(c, "数据集绑定成功")
}
//获取已绑定的数据集
func (es *EvaluationsetService) GetBindedDataset(c *gin.Context){
	id := c.Param("id")
	if id == "" {
		Error(c, errors.New("缺少评测集ID参数"))
		return
	}
	mappings, err := edmc.GetMappingsByEvaluationsetID(id)
	if err != nil {
		Error(c, err)
		return
	}
	datasetIDs := make([]string, 0, len(mappings))
	for _, mapping := range mappings {
		datasetIDs = append(datasetIDs, mapping.DatasetID)
	}

	// 根据数据集ID批量获取数据集详情
	datasets, err := dc.GetDatasetByBatch(datasetIDs)
	if err != nil {
		Error(c, err)
		return
	}

	// 创建映射ID到数据集的关系
	mappingMap := make(map[string]string)
	for _, mapping := range mappings {
		mappingMap[mapping.DatasetID] = mapping.EvaluationsetDatasetMappingID
	}

	// 构建包含映射ID的结果
	result := make([]model.DatasetWithMapping, 0, len(datasets))
	for _, dataset := range datasets {
		mappingID := mappingMap[dataset.DatasetID]
		datasetWithMapping := model.DatasetWithMapping{
			MappingID:   mappingID,
			DatasetID:   dataset.DatasetID,
			Name:        dataset.Name,
			DataCount:   dataset.DataCount,
			Description: dataset.Description,
			CreatedAt:   dataset.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		result = append(result, datasetWithMapping)
	}
	
	Success(c, result)
}
//解除绑定数据集
func (es *EvaluationsetService) UnbindDataset(c *gin.Context){
	id := c.Param("id")
	if id == "" {
		Error(c, errors.New("缺少数据集ID参数"))
		return
	}
	err := edmc.DeleteDatasetMapping(id)
	if err != nil {
		Error(c, err)
		return
	}

	Success(c, "数据集解绑成功")
}
//解除绑定数据集批量处理
func (es *EvaluationsetService) UnbindDatasetBatch(c *gin.Context){
	var mappingIDs []string
	if err := c.ShouldBindJSON(&mappingIDs); err != nil {
		Error(c, err)
		return
	}
	err := edmc.DeleteDatasetMappingByBatch(mappingIDs)
	if err != nil {
		Error(c, err)
		return
	}

	Success(c, "数据集解绑成功")
}


func (es *EvaluationsetService) DeleteEvaluationset(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		Error(c, errors.New("缺少评测集ID参数"))
		return
	}

	ec := &sql.EvaluationsetCRUD{}
	err := ec.DeleteEvaluationset(id)
	if err != nil {
		Error(c, err)
		return
	}

	Success(c, gin.H{"message": "评测集删除成功"})
}

func EditEvaluationset(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func SearchEvaluationsets(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func SortEvaluationsetsByNameUp(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func SortEvaluationsetsByNameDown(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func SortEvaluationsetsByTimeUp(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func SortEvaluationsetsByTimeDown(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
