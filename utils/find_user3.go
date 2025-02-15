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

// FindUser3 ищет пользователя 3 в коллекции users3 по имени пользователя
func FindUser3(username string) (models.User3, error) {
	var user models.User3

	// Получаем коллекцию users3
	collection := database.DB.Collection("users3")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Поиск пользователя по полю username
	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		log.Println("Ошибка при поиске пользователя:", err)
		if err.Error() == "mongo: no documents in result" {
			return user, errors.New("Пользователь не найден")
		}
		return user, err
	}

	return user, nil
}
