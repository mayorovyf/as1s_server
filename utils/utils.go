// utils/utils.go
package utils

import (
	"crypto/rand"
	"encoding/base64"
	"log"
)

// GenerateAPIKey генерирует уникальный API ключ
func GenerateAPIKey() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		log.Println("Ошибка генерации API ключа:", err)
		return "", err
	}
	apiKey := base64.URLEncoding.EncodeToString(key)
	return apiKey, nil
}
