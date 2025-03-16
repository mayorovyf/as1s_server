// utils/insert_user3.go
package utils

import (
	"as1s_server/database"
	"as1s_server/models"
	"context"
	"log"
	"time"
)

// InsertUser3 вставляет нового пользователя в коллекцию users3
func InsertUser3(user models.User3) error {
	// Получаем коллекцию users3 из базы данных
	collection := database.DB.Collection("users3")

	// Устанавливаем тайм-аут для операции
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Вставляем нового пользователя в коллекцию
	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		log.Println("Ошибка при вставке пользователя в базу данных:", err)
		return err
	}

	return nil
}
