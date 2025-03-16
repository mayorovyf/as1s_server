// main.go
package main

import (
	"as1s_server/config"
	"as1s_server/database"
	"as1s_server/handlers"
	"as1s_server/middleware"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	config.LoadEnv()
	database.ConnectDB()

	r := gin.Default()

	// Роуты регистрации
	r.POST("/register_user1", handlers.RegisterUser1)
	r.POST("/register_user2", handlers.RegisterUser2)
	r.POST("/register_user3", handlers.RegisterUser3)

	// Роуты логина
	r.POST("/login_user1", handlers.LoginUser1)
	r.POST("/login_user2", handlers.LoginUser2)
	r.POST("/login_user3", handlers.LoginUser3)

	// Роуты получения данных с использованием middleware
	r.GET("/user1_data", middleware.AuthUser1(), handlers.GetUser1Data)
	r.GET("/user2_data", middleware.AuthUser2(), handlers.GetUser2Data)
	r.GET("/user3_data", middleware.AuthUser3(), handlers.GetUser3Data)

	// Прочие эндпоинты
	r.POST("/verify_qr_user3", middleware.AuthUser3(), handlers.VerifyQRForUser3)
	r.POST("/get_class_users", middleware.AuthUser2(), handlers.GetClassUsersForUser2)
	r.POST("/add_student", middleware.AuthUser2(), handlers.AddStudentToClass)
	r.POST("/update_qr", middleware.AuthUser2(), handlers.UpdateStudentQR)

	// Запуск сервера
	port := config.GetEnv("PORT", "8080")
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
