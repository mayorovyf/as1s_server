// handlers/login_user3_route.go
package handlers

import (
	"as1s_server/models"
	"as1s_server/utils"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// LoginUser3 обрабатывает запрос на логин пользователя 3
func LoginUser3(c *gin.Context) {
	var loginReq models.LoginRequest

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	log.Printf("Попытка входа пользователя: %s", loginReq.Username)

	user, err := utils.FindUser3(loginReq.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не найден"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный пароль"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"api_key": user.APIKey})
}
