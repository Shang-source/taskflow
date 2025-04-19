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

// GenerateToken 生成 JWT
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

// AuthRequired 作为 Gin 中间件，验证 JWT
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
				return nil, fmt.Errorf("不支持的签名方法")
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token无效或过期"})
			return
		}

		// 从 token 中获取 claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := int64(claims["userID"].(float64))
			username := claims["username"].(string)

			// 将 userID, username 存储到上下文中，后续处理可以拿来用
			c.Set("userID", userID)
			c.Set("username", username)
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "解析Token出错"})
			return
		}

		c.Next()
	}
}
