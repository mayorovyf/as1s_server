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

// VerifyQRRequest представляет тело запроса для верификации QR кода от user3.
type VerifyQRRequest struct {
	APIKey string `json:"api_key"`
	QRData string `json:"qr_data"`
}

// VerifyQRForUser3 обрабатывает запрос от user3 на подтверждение QR кода.
// Принимает API ключ user3 и данные QR. Проверяет API ключ, ищет user1 по QR,
// и если у найденного user1 поле used == false, обновляет его на true.
// В ответе возвращается статус true при успехе или false при ошибке.
func VerifyQRForUser3(c *gin.Context) {
	var req VerifyQRRequest

	// Декодирование JSON из тела запроса
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[ERROR] Неверный формат запроса: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  "Неверный формат запроса",
		})
		return
	}

	// Проверка обязательных полей
	if req.APIKey == "" || req.QRData == "" {
		log.Printf("[ERROR] Отсутствуют обязательные поля: APIKey='%s', QRData='%s'", req.APIKey, req.QRData)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  "API ключ и qr_data обязательны",
		})
		return
	}

	// Проверка достоверности API ключа пользователя 3
	user3, err := utils.FindUser3ByAPIKey(req.APIKey)
	if err != nil {
		log.Printf("[ERROR] Неверный API ключ '%s': %v", req.APIKey, err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": false,
			"error":  "Неверный API ключ",
		})
		return
	}
	log.Printf("[INFO] User3 '%s' аутентифицирован успешно", user3.Username)

	// Поиск пользователя 1 по qr_data в коллекции users1
	collection := database.DB.Collection("users1")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user1 models.User1
	err = collection.FindOne(ctx, bson.M{"qr": req.QRData}).Decode(&user1)
	if err != nil {
		log.Printf("[ERROR] Пользователь с QR '%s' не найден: %v", req.QRData, err)
		c.JSON(http.StatusNotFound, gin.H{
			"status": false,
			"error":  "Пользователь с таким QR не найден",
		})
		return
	}
	log.Printf("[INFO] Найден user1 '%s' с QR '%s'", user1.Username, req.QRData)

	// Проверка, что QR код ещё не использован
	if user1.Used {
		log.Printf("[WARN] QR код '%s' уже использован для user1 '%s'", req.QRData, user1.Username)
		c.JSON(http.StatusConflict, gin.H{
			"status": false,
			"error":  "QR код уже использован",
		})
		return
	}

	// Обновляем поле used для найденного user1 на true
	update := bson.M{"$set": bson.M{"used": true}}
	updateCtx, updateCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer updateCancel()

	result, err := collection.UpdateOne(updateCtx, bson.M{"qr": req.QRData}, update)
	if err != nil {
		log.Printf("[ERROR] Ошибка при обновлении статуса QR '%s' для user1 '%s': %v", req.QRData, user1.Username, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"error":  "Ошибка обновления статуса QR кода",
		})
		return
	}
	if result.ModifiedCount == 0 {
		log.Printf("[ERROR] Обновление не выполнено для QR '%s'", req.QRData)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"error":  "Не удалось обновить статус QR кода",
		})
		return
	}

	log.Printf("[INFO] QR код '%s' успешно подтвержден. user1 '%s' обновлен, user3 '%s' уведомлен.", req.QRData, user1.Username, user3.Username)
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "QR код подтвержден и статус обновлен.",
		"user1":   user1.Username,
		"user3":   user3.Username,
	})
}
