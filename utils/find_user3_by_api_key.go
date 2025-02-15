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

// FindUser3ByAPIKey ищет пользователя в коллекции users3 по API ключу
func FindUser3ByAPIKey(apiKey string) (models.User3, error) {
	var user models.User3
	collection := database.DB.Collection("users3")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, bson.M{"api_key": apiKey}).Decode(&user)
	if err != nil {
		log.Println("Ошибка при поиске пользователя по API ключу:", err)
		return user, errors.New("Пользователь не найден")
	}
	return user, nil
}
