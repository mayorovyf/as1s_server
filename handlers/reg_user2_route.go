package handlers

import (
	"as1s_server/models"
	"as1s_server/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterUser2 обрабатывает запрос на регистрацию пользователя 2
func RegisterUser2(c *gin.Context) {
	var user models.User2

	// Декодирование JSON из тела запроса
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	// Генерация API ключа
	apiKey, err := utils.GenerateAPIKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации API ключа"})
		return
	}
	user.APIKey = apiKey

	// Сохранение пользователя в коллекцию users2
	if err := utils.InsertUser2(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка регистрации пользователя"})
		return
	}

	// Ответ
	c.JSON(http.StatusOK, gin.H{"message": "User2 registered with API key: " + apiKey})
}
