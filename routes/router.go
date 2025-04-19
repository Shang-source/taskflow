package routes

import (
	"github.com/gin-gonic/gin"
	"taskflow/controllers"
	"taskflow/middlewares"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 可选: 全局日志中间件
	r.Use(middlewares.Logger())

	// 用户相关路由
	r.POST("/register", controllers.RegisterUser)
	r.POST("/login", controllers.LoginUser)

	// 需要登录后才可访问
	auth := r.Group("/")
	auth.Use(middlewares.AuthRequired())
	{
		auth.POST("/tasks", controllers.CreateTask)
		auth.GET("/tasks", controllers.GetTasks)
		auth.GET("/tasks/:id", controllers.GetTaskByID)
		auth.PUT("/tasks/:id", controllers.UpdateTask)
		auth.DELETE("/tasks/:id", controllers.DeleteTask)
		// routes/router.go
		auth.GET("/notifications", controllers.GetNotifications)
	}

	return r
}
