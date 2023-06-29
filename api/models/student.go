package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type Student struct {
	gorm.Model
	Username     string `gorm:"size:100;not null"              json:"username"`
	Email        string `gorm:"type:varchar(100);unique_index" json:"email"`
	Password     string `gorm:"size:100;not null"              json:"password"`
	ProfileImage string `gorm:"size:255"                       json:"profileimage"`
	Books        []Book
	Notification Notification
}

// hashPassword
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	return string(bytes), err
}

// CheckPasswordHash
func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return errors.New("password incorrect")
	}
	return nil
}
