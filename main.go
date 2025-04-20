package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"os"
	"strconv"
	"taskflow/config"
	_ "taskflow/docs" // import swagger document
	"taskflow/routes"
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
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Not fund.env file，using system env")
	}
	// connect database
	config.ConnectDB()
	defer config.DB.Close()

	// get port id
	portStr := os.Getenv("RUN_PORT")
	if portStr == "" {
		portStr = "8080"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = 8080
	}

	// initialize Gin
	r := routes.SetupRouter()

	// register Swagger route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// start service
	log.Printf("Server running at http://localhost:%d\n", port)
	r.Run(fmt.Sprintf(":%d", port))
}
