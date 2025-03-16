// utils/find_user1_by_api_key.go
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

// FindUser1ByAPIKey ищет пользователя в коллекции users1 по API ключу
func FindUser1ByAPIKey(apiKey string) (models.User1, error) {
	var user models.User1
	collection := database.DB.Collection("users1")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, bson.M{"api_key": apiKey}).Decode(&user)
	if err != nil {
		log.Println("Ошибка при поиске пользователя по API ключу:", err)
		return user, errors.New("Пользователь не найден")
	}
	return user, nil
}
