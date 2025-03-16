// handlers/user2_data_route.go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUser2Data возвращает данные пользователя user2, полученные из контекста.
func GetUser2Data(c *gin.Context) {
	user, exists := c.Get("user2")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Пользователь не найден в контексте"})
		return
	}
	c.JSON(http.StatusOK, user)
}
