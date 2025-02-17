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

// GetClassUsersRequest представляет входные данные запроса для получения списка учеников по id класса.
type GetClassUsersRequest struct {
	APIKey string `json:"api_key"`
}

// GetClassUsersForUser2 проверяет API ключ пользователя 2, получает его class_id и возвращает список всех user1 с таким же class_id.
func GetClassUsersForUser2(c *gin.Context) {
	var req GetClassUsersRequest

	// Декодирование JSON из тела запроса
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[ERROR] Неверный формат запроса: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    false,
			"error":     "Неверный формат запроса",
			"timestamp": time.Now().Unix(),
		})
		return
	}

	// Проверка обязательного поля APIKey
	if req.APIKey == "" {
		log.Printf("[ERROR] API ключ отсутствует в запросе")
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    false,
			"error":     "API ключ обязателен",
			"timestamp": time.Now().Unix(),
		})
		return
	}

	// Проверка действительности API ключа для пользователя 2
	user2, err := utils.FindUser2ByAPIKey(req.APIKey)
	if err != nil {
		log.Printf("[ERROR] Пользователь 2 не найден для API ключа '%s': %v", req.APIKey, err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":    false,
			"error":     "Неверный API ключ",
			"timestamp": time.Now().Unix(),
		})
		return
	}
	log.Printf("[INFO] User2 '%s' аутентифицирован успешно", user2.Username)

	// Извлечение id класса из user2
	classID := user2.ClassID
	if classID == "" {
		log.Printf("[ERROR] У пользователя 2 '%s' отсутствует id класса", user2.Username)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":    false,
			"error":     "Не удалось получить id класса для пользователя 2",
			"timestamp": time.Now().Unix(),
		})
		return
	}
	log.Printf("[INFO] Для user2 '%s' получен class_id: %s", user2.Username, classID)

	// Поиск всех user1 с таким же class_id в коллекции users1
	collection := database.DB.Collection("users1")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{"class_id": classID})
	if err != nil {
		log.Printf("[ERROR] Ошибка при поиске пользователей с class_id '%s': %v", classID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":    false,
			"error":     "Ошибка при поиске пользователей по id класса",
			"timestamp": time.Now().Unix(),
		})
		return
	}
	defer cursor.Close(ctx)

	var users []models.User1
	for cursor.Next(ctx) {
		var user models.User1
		if err := cursor.Decode(&user); err != nil {
			log.Printf("[ERROR] Ошибка декодирования документа пользователя: %v", err)
			continue
		}
		// Убираем поле password из результата
		user.Password = ""
		users = append(users, user)
	}
	if err := cursor.Err(); err != nil {
		log.Printf("[ERROR] Ошибка работы курсора: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":    false,
			"error":     "Ошибка обработки результатов поиска",
			"timestamp": time.Now().Unix(),
		})
		return
	}

	// Если не найдено ни одного пользователя
	if len(users) == 0 {
		log.Printf("[INFO] Пользователи с class_id '%s' не найдены", classID)
		c.JSON(http.StatusNotFound, gin.H{
			"status":    false,
			"error":     "Ученики с указанным id класса не найдены",
			"timestamp": time.Now().Unix(),
		})
		return
	}

	// Успешный ответ
	log.Printf("[INFO] Найдено %d пользователей с class_id '%s'", len(users), classID)
	c.JSON(http.StatusOK, gin.H{
		"status":    true,
		"users":     users,
		"timestamp": time.Now().Unix(),
	})
}
