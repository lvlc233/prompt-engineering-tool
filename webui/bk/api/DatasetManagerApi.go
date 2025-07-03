package api

import (
	"webui-server/service"

	"github.com/gin-gonic/gin"
)

var ds *service.DatasetService

func SetDatasetManagerRouter(r *gin.Engine) {

	DatasetMapManagerRouter := r.Group("/dataset_manager")
	{
		//获取所有数据集
		DatasetMapManagerRouter.GET("/", ds.GetDatasets)
		//获取单个数据集基本信息 
		DatasetMapManagerRouter.GET("/info/:id", ds.GetDatasetInfo)
		//获取数据集详细数据
		DatasetMapManagerRouter.GET("/:id", ds.GetDatasetDetail)

		//添加数据集
		DatasetMapManagerRouter.POST("/add", ds.CreateDataset)

		//查看详情
		//删除
		DatasetMapManagerRouter.DELETE("/:id", ds.DeleteDataset)
		//编辑
		DatasetMapManagerRouter.POST("/editor", ds.EditDataset)

	}
}
