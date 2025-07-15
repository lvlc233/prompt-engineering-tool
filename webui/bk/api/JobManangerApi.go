package api

import (
	"webui-server/service"

	"github.com/gin-gonic/gin"
)
var js = &service.JobService{}


func SetJobMapManagerRouter(r *gin.Engine) {

	JobMapManagerRouter := r.Group("/job_manager")
	{
		//获取所有任务
		JobMapManagerRouter.GET("/", js.GetJobs)
		//获取任务的所有版本
		JobMapManagerRouter.GET("/versions/:id", js.GetJobVersionsByID)
		//创建任务
		JobMapManagerRouter.POST("/add", js.CreateJob)
		//删除任务
		JobMapManagerRouter.DELETE("/:id", js.DeleteJob)
	
		//编辑	
		//设置评价标准
		JobMapManagerRouter.POST("/set/criteria",es.SetEvaluationCriteria)
		//设置分数上限
		JobMapManagerRouter.POST("/set/score_cap",es.SetScoreCap)
		
	
		//绑定数据集合
		JobMapManagerRouter.POST("/bind_dataset", es.BindDataset)
		//获取绑定的数据集
		JobMapManagerRouter.GET("/bind_dataset/:id", es.GetBindedDataset)
		//解绑数据集
		JobMapManagerRouter.DELETE("/unbind_dataset/:id", es.UnbindDataset)
		//解绑数据集批次
		JobMapManagerRouter.DELETE("/unbind_dataset_batch", es.UnbindDatasetBatch)
	}
}