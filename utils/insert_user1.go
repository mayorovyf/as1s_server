package utils

import (
	"as1s_server/database"
	"as1s_server/models"
	"context"
	"log"
	"time"
)

// InsertUser вставляет нового пользователя в коллекцию users1
func InsertUser1(user models.User1) error {
	collection := database.DB.Collection("users1")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		log.Println("Ошибка при вставке пользователя в базу данных:", err)
		return err
	}
	return nil
}
