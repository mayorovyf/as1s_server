// models/user2.go
package models

// User2 структура
type User2 struct {
	FirstName  string `json:"first_name" bson:"first_name"`
	MiddleName string `json:"middle_name" bson:"middle_name"`
	LastName   string `json:"last_name" bson:"last_name"`
	Username   string `json:"username" bson:"username"`
	Password   string `json:"password" bson:"password"`
	Class      string `json:"class" bson:"class"`
	ClassID    string `json:"class_id" bson:"class_id"`
	QR         string `json:"qr" bson:"qr"`
	Used       bool   `json:"used" bson:"used"`
	InBuilding bool   `json:"in_building" bson:"in_building"`
	APIKey     string `json:"api_key" bson:"api_key"`
}
