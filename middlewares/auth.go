package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// GenerateToken JWT
func GenerateToken(userID int64, username string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	expireHoursStr := os.Getenv("JWT_EXPIRE_HOURS")
	if secret == "" || expireHoursStr == "" {
		return "", fmt.Errorf("JWT 配置缺失")
	}

	expireHours, err := strconv.Atoi(expireHoursStr)
	if err != nil {
		expireHours = 24
	}

	claims := jwt.MapClaims{
		"userID":   userID,
		"username": username,
		"exp":      time.Now().Add(time.Duration(expireHours) * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// AuthRequired middleware for Gin to validate JWT
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "缺少Token"})
			return
		}

		tokenString := authHeader
		secret := os.Getenv("JWT_SECRET")

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unsupported signing method")
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Extract claims from the token
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := int64(claims["userID"].(float64))
			username := claims["username"].(string)

			// Store userID and username in the context for later use in subsequent handlers.
			c.Set("userID", userID)
			c.Set("username", username)
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "解析Token出错"})
			return
		}

		c.Next()
	}
}
