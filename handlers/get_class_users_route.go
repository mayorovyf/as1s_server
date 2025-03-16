// handlers/get_class_users_route.go
package handlers

import (
	"as1s_server/database"
	"as1s_server/models"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// GetClassUsersForUser2 возвращает список учеников, основываясь на class_id учителя (user2), извлечённого из контекста.
func GetClassUsersForUser2(c *gin.Context) {
	// Извлекаем учителя из контекста
	teacherInterface, exists := c.Get("user2")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": false, "error": "Учитель не аутентифицирован"})
		return
	}
	teacher, ok := teacherInterface.(models.User2)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "error": "Ошибка получения данных учителя"})
		return
	}

	classID := teacher.ClassID
	if classID == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "error": "Не удалось получить id класса для учителя"})
		return
	}

	collection := database.DB.Collection("users1")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{"class_id": classID})
	if err != nil {
		log.Printf("[ERROR] Ошибка поиска пользователей с class_id '%s': %v", classID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "error": "Ошибка поиска учеников"})
		return
	}
	defer cursor.Close(ctx)

	var users []models.User1
	for cursor.Next(ctx) {
		var user models.User1
		if err := cursor.Decode(&user); err != nil {
			log.Printf("[ERROR] Ошибка декодирования пользователя: %v", err)
			continue
		}
		// Очистка поля пароля для безопасности
		user.Password = ""
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("[ERROR] Ошибка работы курсора: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "error": "Ошибка обработки результатов поиска"})
		return
	}
	if len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": false, "error": "Ученики с указанным id класса не найдены"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": true, "users": users})
}
