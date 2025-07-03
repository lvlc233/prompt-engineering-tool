package service

import (
	"github.com/gin-gonic/gin"
)

func GetEvaluationsets(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func GetEvaluationsetDetail(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func CreateEvaluationset(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func DeleteEvaluationset(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
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
