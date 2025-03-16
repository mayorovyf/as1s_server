// config/config.go
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv загружает переменные окружения из .env
func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден, используются стандартные переменные окружения")
	}
}

// GetEnv возвращает значение переменной окружения
func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
