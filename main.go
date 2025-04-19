package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"taskflow/config"
	_ "taskflow/docs" // 导入 swagger 文档
	"taskflow/jobs"
	"taskflow/routes"

	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title TaskFlow API
// @version 1.0
// @description Swagger for Gin-based TaskFlow app
// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// 加载环境变量
	err := godotenv.Load()
	if err != nil {
		log.Println("Not fund.env file，using system env")
	}

	// 连接数据库
	config.ConnectDB()
	defer config.DB.Close()

	// 获取端口号
	portStr := os.Getenv("RUN_PORT")
	if portStr == "" {
		portStr = "8080"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = 8080
	}

	// 初始化 Gin
	r := routes.SetupRouter()

	// 注册 Swagger 路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 启动到期提醒调度
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	jobs.StartScheduler(ctx)

	// 启动服务
	log.Printf("Server running at http://localhost:%d\n", port)
	r.Run(fmt.Sprintf(":%d", port))
}
