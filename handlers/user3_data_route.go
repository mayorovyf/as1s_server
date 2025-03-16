// handlers/user3_data_route.go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUser3Data возвращает данные пользователя user3, полученные из контекста.
func GetUser3Data(c *gin.Context) {
	user, exists := c.Get("user3")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Пользователь не найден в контексте"})
		return
	}
	c.JSON(http.StatusOK, user)
}
