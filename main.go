package main

import (
	"as1s_server/config"
	"as1s_server/database"
	"as1s_server/handlers"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	config.LoadEnv()
	database.ConnectDB()

	r := gin.Default()

	// Роуты
	r.POST("/register_user1", handlers.RegisterUser1)
	r.POST("/register_user2", handlers.RegisterUser2)
	r.POST("/register_user3", handlers.RegisterUser3)
	r.POST("/login_user1", handlers.LoginUser1)
	r.POST("/login_user2", handlers.LoginUser2)
	r.POST("/login_user3", handlers.LoginUser3)
	r.GET("/user1_data", handlers.GetUser1Data)
	r.GET("/user2_data", handlers.GetUser2Data)
	r.GET("/user3_data", handlers.GetUser3Data)
	r.POST("/verify_qr_user3", handlers.VerifyQRForUser3)
	r.POST("/get_class_users", handlers.GetClassUsersForUser2)

	// Запуск сервера
	port := config.GetEnv("PORT", "8080")
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
