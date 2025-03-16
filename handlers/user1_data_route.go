// handlers/user1_data_route.go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUser1Data возвращает данные пользователя user1, полученные из контекста.
func GetUser1Data(c *gin.Context) {
	user, exists := c.Get("user1")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Пользователь не найден в контексте"})
		return
	}
	c.JSON(http.StatusOK, user)
}
