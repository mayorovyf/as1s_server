package models

// LoginRequest структура для запроса при логине
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
