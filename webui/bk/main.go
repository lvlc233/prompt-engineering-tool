package main

import (
	"log"
	"webui-server/api"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	db, err := InitDatabase()
	if err != nil {
		log.Fatal("数据库初始化失败:", err)
	}
	defer db.Close()

	// 创建 Gin 路由
	r := gin.Default()

	// 配置 CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	api.SetDatasetManagerRouter(r)
	api.SetEvaluationsetManagerRouter(r)
	api.SetJobMapManagerRouter(r)
	api.SetSettingRouter(r)

	log.Println("服务器启动在 :8593")
	r.Run(":8593")
}
