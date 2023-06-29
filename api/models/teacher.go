package models

import "github.com/jinzhu/gorm"

type Teacher struct {
	gorm.Model
	Username string `gorm:"size:50;unique_index;not null"    json:"username"`
	Email    string `gorm:"size:100;unique_index;not null"   json:"email"`
	Password string `gorm:"size:100;not null"                json:"password"`
	Notification []Notification
}
