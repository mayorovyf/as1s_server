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

// VerifyQRRequest определяет структуру входных данных для проверки QR кода.
type VerifyQRRequest struct {
	QRData string `json:"qr_data"`
}

// VerifyQRForUser3 подтверждает QR код, используя информацию из контекста user3 (охранника).
// Сначала ищется пользователь user1 по полю qr, затем проверяется, не был ли код использован,
// после чего обновляется документ по уникальному полю username. При этом:
// • В массиве times сохраняются только 5 последних заходов,
// • Поле in_building меняется на противоположное (toggle).
func VerifyQRForUser3(c *gin.Context) {
	var req VerifyQRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[ERROR] Неверный формат запроса: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": "Неверный формат запроса"})
		return
	}
	if req.QRData == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": "Параметр qr_data обязателен"})
		return
	}

	// Извлекаем user3 из контекста.
	user3Interface, exists := c.Get("user3")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": false, "error": "Пользователь не аутентифицирован"})
		return
	}
	user3, ok := user3Interface.(models.User3)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "error": "Ошибка получения данных user3"})
		return
	}

	// Ищем пользователя user1 по QR коду.
	collection := database.DB.Collection("users1")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user1 models.User1
	err := collection.FindOne(ctx, bson.M{"qr": req.QRData}).Decode(&user1)
	if err != nil {
		log.Printf("[ERROR] Пользователь с QR '%s' не найден: %v", req.QRData, err)
		c.JSON(http.StatusNotFound, gin.H{"status": false, "error": "Пользователь с таким QR не найден"})
		return
	}

	if user1.Used {
		log.Printf("[WARN] QR код '%s' уже использован для user1 '%s'", req.QRData, user1.Username)
		c.JSON(http.StatusConflict, gin.H{"status": false, "error": "QR код уже использован"})
		return
	}

	// Если поле times имеет значение null, инициализируем его пустым массивом.
	if user1.Times == nil {
		setTimes := bson.M{"$set": bson.M{"times": []time.Time{}}}
		_, err := collection.UpdateOne(ctx, bson.M{"username": user1.Username}, setTimes)
		if err != nil {
			log.Printf("[WARN] Не удалось инициализировать поле times для user1 '%s': %v", user1.Username, err)
		} else {
			user1.Times = []time.Time{}
		}
	}

	currentTime := time.Now()
	// Вычисляем новое значение поля in_building – оно должно стать противоположным текущему.
	newInBuilding := !user1.InBuilding

	// Обновляем документ:
	// – отмечаем QR как использованный,
	// – фиксируем время использования,
	// – добавляем текущее время в массив times (оставляем только 5 последних),
	// – переключаем значение in_building.
	update := bson.M{
		"$set": bson.M{
			"used":         true,
			"last_used_at": currentTime,
			"in_building":  newInBuilding,
		},
		"$push": bson.M{
			"times": bson.M{
				"$each":  []time.Time{currentTime},
				"$slice": -5,
			},
		},
	}

	updateCtx, updateCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer updateCancel()
	result, err := collection.UpdateOne(updateCtx, bson.M{"username": user1.Username}, update)
	if err != nil || result.ModifiedCount == 0 {
		log.Printf("[ERROR] Не удалось обновить статус QR для user1 '%s': %v", user1.Username, err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "error": "Не удалось обновить статус QR кода"})
		return
	}

	log.Printf("[INFO] QR код '%s' успешно подтвержден. user1: '%s', user3: '%s'. Новое значение in_building: %v", req.QRData, user1.Username, user3.Username, newInBuilding)
	c.JSON(http.StatusOK, gin.H{
		"status":      true,
		"message":     "QR код подтвержден и статус обновлен.",
		"user1":       user1.Username,
		"user3":       user3.Username,
		"in_building": newInBuilding,
	})
}
