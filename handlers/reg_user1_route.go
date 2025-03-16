// handlers/reg_user1_route.go
package handlers

import (
	"as1s_server/models"
	"as1s_server/utils"
	"golang.org/x/crypto/bcrypt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterUser1(c *gin.Context) {
	var user models.User1

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	// Проверка уникальности логина (username)
	if existingUser, err := utils.FindUser(user.Username); err == nil && existingUser.Username != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Пользователь с таким логином уже существует"})
		return
	}

	// Хэширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка хэширования пароля"})
		return
	}
	user.Password = string(hashedPassword)

	// Генерация API ключа
	apiKey, err := utils.GenerateAPIKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации API ключа"})
		return
	}
	user.APIKey = apiKey

	if err := utils.InsertUser1(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка регистрации пользователя"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User1 registered with API key: " + apiKey})
}
