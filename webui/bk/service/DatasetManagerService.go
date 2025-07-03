package service

import (
	"errors"
	"fmt"
	"webui-server/model"
	"webui-server/sql"

	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

type DatasetService struct{}

var dc *sql.DatasetCRUD
var ddc *sql.DatasetDetailCRUD

// 获取所有的数据集
func (ds *DatasetService) GetDatasets(c *gin.Context) {

	datasets, err := dc.GetAllDatasets()
	if err != nil {
		Error(c, err)
		return
	}
	// 返回成功响应
	Success(c, datasets)

}

// 创建数据集
func (ds *DatasetService) CreateDataset(c *gin.Context) {

	var createDatasetRequest model.CreateDatasetRequest
	//参数解析
	if err := c.ShouldBindJSON(&createDatasetRequest); err != nil {
		Error(c, err)
		return
	}
	if err := dc.AddDataset(createDatasetRequest.Name, createDatasetRequest.Description); err != nil {
		Error(c, err)
		return
	}
	Success(c, "数据集创建成功")

}

// 获取单个数据集基本信息
func (ds *DatasetService) GetDatasetInfo(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		Error(c, errors.New("缺少数据集ID参数"))
		return
	}
	dataset, err := dc.GetDatasetByID(id)
	if err != nil {
		Error(c, err)
		return
	}
	Success(c, dataset)
}

// 获取数据集具体数据
func (ds *DatasetService) GetDatasetDetail(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		Error(c, errors.New("缺少数据集ID参数"))
		return
	}
	dataset_detail, err := ddc.GetDatasetDetailByDatasetID(id)
	if err != nil {
		Error(c, err)
		return
	}
	Success(c, dataset_detail)
}

// 删除数据集
func (ds *DatasetService) DeleteDataset(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		Error(c, errors.New("缺少数据集ID参数"))
		return
	}
	if err := dc.DeleteDataset(id); err != nil {
		Error(c, err)
		return
	}
	Success(c, "数据集删除成功")
}

// 包括增/删改
// 逻辑如下
// 接受map表示操作
// 若id为空,表示新增
// 若id不为空,且input/output有数据为修改
// 若id不为空,而input/output皆无数据为删除
func (ds *DatasetService) EditDataset(c *gin.Context) {

	var editDatasetRequest model.EditDatasetRequest
	//参数解析
	if err := c.ShouldBindJSON(&editDatasetRequest); err != nil {
		Error(c, err)
		return
	}
	datasetID := editDatasetRequest.ID
	name := editDatasetRequest.Name
	desc := editDatasetRequest.Description
	tuples := editDatasetRequest.EditDataTuples

	fmt.Println(tuples)

	// 处理数据集基本信息的修改
	if name != "" || desc != "" {
		if err := dc.UpdateDataset(datasetID, name, desc); err != nil {
			Error(c, err)
			return
		}
	}

	// 处理数据元组的增删改操作
	if len(tuples) > 0 {
		// 分类处理不同操作
		var toCreate []model.EditDataTuple // 新增列表
		var toUpdate []model.EditDataTuple
		var toDelete []string // 删除列表（存储ID）

		for _, tuple := range tuples {
			if tuple.ID == "" {
				// ID为0，表示新增
				toCreate = append(toCreate, tuple)
			} else if tuple.Input == "" && tuple.Output == "" {
				// ID不为0且input/output皆无数据，表示删除
				toDelete = append(toDelete, tuple.ID)
			} else {
				// ID不为0且有数据，表示修改
				toUpdate = append(toUpdate, tuple)
			}
		}

		// 执行新增操作
		if len(toCreate) > 0 {
			if err := ddc.AddDatasetDetailByBatch(datasetID, toCreate); err != nil {
				Error(c, err)
				return
			}
		}

		// 执行删除操作
		if len(toDelete) > 0 {
			if err := ddc.DeleteDatasetDetailByBatch(datasetID, toDelete); err != nil {
				Error(c, err)
				return
			}
		}

		// 执行修改操作（批量处理）
		if len(toUpdate) > 0 {
			if err := ddc.UpdateDatasetDetailByBatch(datasetID, toUpdate); err != nil {
				Error(c, err)
				return
			}
		}

	}

	Success(c, "数据集编辑成功")
}
