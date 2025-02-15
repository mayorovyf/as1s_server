package handlers

import (
	"as1s_server/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUser1Data возвращает данные пользователя из коллекции users1 по API ключу
func GetUser1Data(c *gin.Context) {
	apiKey := c.Query("api_key")
	if apiKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "API ключ обязателен"})
		return
	}

	user, err := utils.FindUser1ByAPIKey(apiKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не найден"})
		return
	}

	// Удаляем поле password перед отправкой ответа
	user.Password = ""
	c.JSON(http.StatusOK, user)
}
