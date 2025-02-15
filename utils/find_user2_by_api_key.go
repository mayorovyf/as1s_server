package utils

import (
	"as1s_server/database"
	"as1s_server/models"
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// FindUser2ByAPIKey ищет пользователя в коллекции users2 по API ключу
func FindUser2ByAPIKey(apiKey string) (models.User2, error) {
	var user models.User2
	collection := database.DB.Collection("users2")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, bson.M{"api_key": apiKey}).Decode(&user)
	if err != nil {
		log.Println("Ошибка при поиске пользователя по API ключу:", err)
		return user, errors.New("Пользователь не найден")
	}
	return user, nil
}
