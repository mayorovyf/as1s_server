package handlers

import (
	"as1s_server/models"
	"as1s_server/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterUser1(c *gin.Context) {
	var user models.User1

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

	// Сохранение пользователя в коллекцию
	if err := utils.InsertUser1(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка регистрации пользователя"})
		return
	}

	// Ответ
	c.JSON(http.StatusOK, gin.H{"message": "User1 registered with API key: " + apiKey})
}
