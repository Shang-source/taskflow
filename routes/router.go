package routes

import (
	"github.com/gin-gonic/gin"
	"taskflow/controllers"
	"taskflow/middlewares"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	//  Global logging middleware
	r.Use(middlewares.Logger())

	// related users route
	r.POST("/register", controllers.RegisterUser)
	r.POST("/login", controllers.LoginUser)

	//Authentication required to access
	auth := r.Group("/")
	auth.Use(middlewares.AuthRequired())
	{
		auth.GET("/tasks", controllers.GetTasks)
		// routes/router.go
	}

	return r
}
