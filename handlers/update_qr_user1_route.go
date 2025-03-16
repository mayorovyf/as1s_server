// handlers/update_qr_user1_route.go
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

// UpdateQRRequest – структура входных данных для обновления QR кода (только username ученика).
type UpdateQRRequest struct {
	Username string `json:"username"` // Логин ученика, для которого обновляется QR код
}

// UpdateStudentQR обновляет QR код ученика, используя данные учителя (user2) из контекста.
func UpdateStudentQR(c *gin.Context) {
	var req UpdateQRRequest
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
		c.JSON(http.StatusNotFound, gin.H{"status": false, "error": "Ученик не найден"})
		return
	}

	// Проверка, что ученик принадлежит классу учителя
	if student.ClassID != teacher.ClassID {
		c.JSON(http.StatusForbidden, gin.H{"status": false, "error": "Ученик не принадлежит вашему классу"})
		return
	}

	collection := database.DB.Collection("users1")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Генерация нового уникального QR кода с ограничением по числу попыток
	var newQR string
	var tries int
	for {
		newQR, err = utils.GenerateAPIKey()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "error": "Ошибка генерации нового QR кода"})
			return
		}

		// Проверяем уникальность QR кода в коллекции
		count, err := collection.CountDocuments(ctx, bson.M{"qr": newQR})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "error": "Ошибка проверки уникальности QR кода"})
			return
		}
		if count == 0 {
			break // Уникальный QR найден
		}
		tries++
		if tries > 5 {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "error": "Не удалось сгенерировать уникальный QR код, попробуйте позже"})
			return
		}
	}

	// Обновление данных ученика: новый QR, сброс флага used и времени использования
	update := bson.M{
		"$set": bson.M{
			"qr":           newQR,
			"used":         false,
			"last_used_at": time.Time{},
		},
	}
	_, err = collection.UpdateOne(ctx, bson.M{"username": student.Username}, update)
	if err != nil {
		log.Printf("[ERROR] Ошибка обновления QR кода для ученика '%s': %v", student.Username, err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "error": "Ошибка обновления QR кода"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true, "message": "QR код успешно обновлен", "new_qr": newQR})
}
