package middlewares

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		end := time.Now()
		log.Printf("request: %s %s, Duration: %v", c.Request.Method, c.Request.URL.Path, end.Sub(start))
	}
}
