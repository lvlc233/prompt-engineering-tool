package api

import (
	"webui-server/service"

	"github.com/gin-gonic/gin"
)

func SetDatasetManagerRouter(r *gin.Engine) {

	DatasetMapManagerRouter := r.Group("/dataset_manager")
	{
		//获取所有数据集
		DatasetMapManagerRouter.GET("/", service.GetDatasets)
		//获取单个数据集基本信息
		DatasetMapManagerRouter.GET("/info/:id", service.GetDatasetInfo)
		//获取数据集详细数据
		DatasetMapManagerRouter.GET("/:id", service.GetDatasetDetail)

		//添加数据集
		DatasetMapManagerRouter.POST("/add", service.CreateDataset)

		//查看详情
		//删除
		DatasetMapManagerRouter.DELETE("/:id", service.DeleteDataset)
		//编辑
		DatasetMapManagerRouter.POST("/editor", service.EditDataset)

	}
}
