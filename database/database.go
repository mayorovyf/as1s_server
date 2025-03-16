// database/database.go
package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectDB() {
	// Получаем логин и пароль из переменных окружения
	username := os.Getenv("MONGO_USERNAME")
	password := os.Getenv("MONGO_PASSWORD")
	uri := os.Getenv("MONGO_URI")

	// Формируем строку подключения с логином и паролем
	connectionString := "mongodb://" + username + ":" + password + "@" + uri

	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal("Ошибка создания клиента MongoDB:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal("Ошибка подключения к MongoDB:", err)
	}

	DB = client.Database(os.Getenv("DB_NAME"))
	log.Println("Успешное подключение к MongoDB")
}
