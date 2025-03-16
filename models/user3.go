// models/user3.go
package models

// User3 структура
type User3 struct {
	FirstName  string `json:"first_name" bson:"first_name"`
	MiddleName string `json:"middle_name" bson:"middle_name"`
	LastName   string `json:"last_name" bson:"last_name"`
	Username   string `json:"username" bson:"username"`
	Password   string `json:"password" bson:"password"`
	APIKey     string `json:"api_key" bson:"api_key"`
}
