package models

import (
	"errors"

	"strings"

	"github.com/badoux/checkmail"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	Username     string `gorm:"size:100;not null"              json:"username"`
	PhoneNumber  int    `gorm:"size:20;not null"               json:"phonenumber"`
	Email        string `gorm:"type:varchar(100);unique_index" json:"email"`
	Password     string `gorm:"size:100;not null"              json:"password"`
	ProfileImage string `gorm:"size:255"                       json:"profileimage"`
	Books        []Book
	Notification Notification
}

// hashPassword hash user input password
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

// BeforeSave hashes student password
func (s *Student) BeforeSave() error {
	password := strings.TrimSpace(s.Password)
	hashedpassword, err := HashPassword(password)
	if err != nil {
		return err
	}
	s.Password = string(hashedpassword)
	return nil
}

// Prepare strips input off white spaces
func (s *Student) Prepare() {
	s.Email = strings.TrimSpace(s.Email)
	s.Username = strings.TrimSpace(s.Username)
	s.ProfileImage = strings.TrimSpace(s.ProfileImage)
}

// Validate user input
func (s *Student) Validate(action string) error {
	switch strings.ToLower(action) {
	case "login":
		if s.Email == "" {
			return errors.New("please enter email address")
		}
		if s.Username == "" {
			return errors.New("please provide username or phone number to login.")
		}
		return nil

	default: //create function where all fields are required
		if s.Username == "" {
			return errors.New("please enter username")
		}
		if s.Email == "" {
			return errors.New("please enter email address")
		}
		if s.Password == "" {
			return errors.New("please provide password ")
		}
		if s.PhoneNumber == 0 {
			return errors.New("please enter phone number")
		}
		if err := checkmail.ValidateFormat(s.Email); err != nil {
			return errors.New("invalid email")
		}
		return nil
	}

}

// func SaveStudent adds student to database
func (s *Student) SaveStudent(db *gorm.DB) (*Student, error) {
	var err error
	err = db.Debug().Create(&s).Error
	if err != nil {
		return &Student{}, err
	}
	return s, nil

}

// Get  Student based on email or phone number
func (s *Student) GetStudent(db *gorm.DB) (*Student, error) {
	account := &Student{}
	if err := db.Debug().Table("students").Where("email= ?", s.Email).Or("phonenumber=?", s.PhoneNumber).First(account).Error; err != nil {
		return nil, err
	}
	return account, nil
}
