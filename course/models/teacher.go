package models

import (
	"gorm.io/gorm"
)

type Teacher struct {
	gorm.Model
	Name    string `json:"name"`
	Gender  string `json:"gender"`
	Age     int    `json:"age"`
	Subject string `json:"subject"`

	IDcard string `json:"id"`
	Phone  string `json:"phone"`
}
