package middlewares

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		end := time.Now()
		log.Printf("请求: %s %s, 耗时: %v", c.Request.Method, c.Request.URL.Path, end.Sub(start))
	}
}
