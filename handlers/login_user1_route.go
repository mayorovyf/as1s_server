// handlers/login_user1_route.go
package handlers

import (
	"as1s_server/models"
	"as1s_server/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// LoginUser1 обрабатывает запрос на логин пользователя
func LoginUser1(c *gin.Context) {
	var loginReq models.LoginRequest

	// Декодирование JSON из тела запроса
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	log.Printf("Попытка входа пользователя: %s", loginReq.Username)

	// Поиск пользователя в базе данных
	user, err := utils.FindUser(loginReq.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не найден"})
		return
	}

	// Сравнение пароля с хэшированным значением
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный пароль"})
		return
	}

	// Возвращаем API ключ (пароль не логируется)
	c.JSON(http.StatusOK, gin.H{"api_key": user.APIKey})
}
