// utils/find_user.go
package utils

import (
	"as1s_server/database"
	"as1s_server/models"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// FindUser находит пользователя по имени пользователя
func FindUser(username string) (models.User1, error) {
	var user models.User1
	collection := database.DB.Collection("users1")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		log.Println("Ошибка при поиске пользователя:", err)
		return user, err
	}
	return user, nil
}
