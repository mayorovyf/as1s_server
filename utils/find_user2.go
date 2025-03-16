// utils/find_user2.go
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

// FindUser2 ищет пользователя 2 в коллекции users2 по имени пользователя
func FindUser2(username string) (models.User2, error) {
	var user models.User2

	// Получаем коллекцию users2
	collection := database.DB.Collection("users2")
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
