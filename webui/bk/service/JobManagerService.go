package service

import (
	"errors"
	"webui-server/model"
	"webui-server/sql"

	"github.com/gin-gonic/gin"
)

type JobService struct{}

var jc = &sql.JobCRUD{}

// 获取所有任务
func (js *JobService) GetJobs(c *gin.Context) {
	jobs, err := jc.GetAllJobs()
	if err != nil {
		Error(c, err)
		return
	}
	Success(c, jobs)
}

//
func (js *JobService) GetJobVersionsByID(c *gin.Context){
	id := c.Param("id")
	versions,err:=jc.GetJobVersionsByID(id)
		if err != nil {
		Error(c, err)
		return
	}
	Success(c, versions)

}

// 获取任务详情根据任务ID
func (js *JobService) GeJobDetail(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		Error(c, errors.New("缺少任务ID参数"))
		return
	}
	jobDetail, err := jc.GetJobDetail(id)
	if err != nil {
		Error(c, err)
		return
	}
	if jobDetail == nil {
		Error(c, errors.New("任务不存在"))
		return
	}
	Success(c, jobDetail)
}

// 创建新任务
func (js *JobService) CreateJob(c *gin.Context) {
	var req model.CreateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, err)
		return
	}

	err := jc.AddJob(req.Name,req.Description)
	if err != nil {
		Error(c, err)
		return
	}

	Success(c, "任务创建成功")
}

// 删除任务
func (js *JobService) DeleteJob(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		Error(c, errors.New("缺少任务ID参数"))
		return
	}

	err := jc.DeleteJob(id)
	if err != nil {
		Error(c, err)
		return
	}

	Success(c, gin.H{"message": "任务删除成功"})
}

