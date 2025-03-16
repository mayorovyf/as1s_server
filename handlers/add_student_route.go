// handlers/add_student_route.go
package handlers

import (
	"as1s_server/database"
	"as1s_server/models"
	"as1s_server/utils"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// AddStudentRequest – структура входных данных для добавления ученика (только username).
type AddStudentRequest struct {
	Username string `json:"username"` // Логин ученика
}

// AddStudentToClass обновляет данные ученика, привязывая его к классу учителя (user2), извлечённого из контекста.
func AddStudentToClass(c *gin.Context) {
	var req AddStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[ERROR] Неверный формат запроса: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": "Неверный формат запроса"})
		return
	}
	if req.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": "Параметр username ученика обязателен"})
		return
	}

	// Извлекаем учителя (user2) из контекста
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

	// Поиск ученика (user1) по username
	student, err := utils.FindUser(req.Username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": false, "error": "Ученик не найден. Возможно, его нужно зарегистрировать"})
		return
	}

	// Обновление данных ученика: привязка к классу учителя
	collection := database.DB.Collection("users1")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"class":    teacher.Class,
			"class_id": teacher.ClassID,
		},
	}
	_, err = collection.UpdateOne(ctx, bson.M{"username": student.Username}, update)
	if err != nil {
		log.Printf("[ERROR] Ошибка обновления данных ученика: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "error": "Ошибка обновления данных ученика"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true, "message": "Ученик успешно добавлен/обновлен в классе", "student": student.Username})
}
