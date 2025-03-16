// middleware/auth.go
package middleware

import (
	"net/http"

	"as1s_server/utils"
	"github.com/gin-gonic/gin"
)

// AuthUser1 проверяет API ключ для пользователя user1.
// Если аутентификация проходит успешно, информация о пользователе сохраняется в контексте.
func AuthUser1() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "API ключ обязателен"})
			c.Abort()
			return
		}
		user, err := utils.FindUser1ByAPIKey(apiKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь user1 не найден"})
			c.Abort()
			return
		}
		c.Set("user1", user)
		c.Next()
	}
}

// AuthUser2 проверяет API ключ для пользователя user2 (учителя).
func AuthUser2() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "API ключ обязателен"})
			c.Abort()
			return
		}
		user, err := utils.FindUser2ByAPIKey(apiKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь user2 не найден"})
			c.Abort()
			return
		}
		c.Set("user2", user)
		c.Next()
	}
}

// AuthUser3 проверяет API ключ для пользователя user3 (охранника).
func AuthUser3() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "API ключ обязателен"})
			c.Abort()
			return
		}
		user, err := utils.FindUser3ByAPIKey(apiKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь user3 не найден"})
			c.Abort()
			return
		}
		c.Set("user3", user)
		c.Next()
	}
}
