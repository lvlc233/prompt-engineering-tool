package service

import (
	"github.com/gin-gonic/gin"
)

func GetJobs(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func GetJobDetail(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func CreateJob(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func DeleteJob(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func CreateJobVersion(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func SwitchJobVersion(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func GetJobVersions(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func SearchJobs(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func SortJobsByNameUp(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func SortJobsByNameDown(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func SortJobsByTimeUp(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func SortJobsByTimeDown(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
