package handlers

import (
	"as1s_server/models"
	"as1s_server/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// LoginUser3 обрабатывает запрос на логин пользователя 3
func LoginUser3(c *gin.Context) {
	var loginReq models.LoginRequest

	// Декодирование JSON из тела запроса
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	// Логирование попытки входа
	log.Printf("Попытка входа пользователя: %s", loginReq.Username)

	// Поиск пользователя в базе данных
	user, err := utils.FindUser3(loginReq.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не найден"})
		return
	}

	// Логирование пароля из запроса и пароля в базе
	log.Printf("Пароль из запроса: %s", loginReq.Password)
	log.Printf("Пароль в базе данных: %s", user.Password)

	// Сравнение пароля
	if loginReq.Password != user.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный пароль"})
		return
	}

	// Ответ с API ключом
	c.JSON(http.StatusOK, gin.H{"api_key": user.APIKey})
}
