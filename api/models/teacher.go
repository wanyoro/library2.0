package models

import (
	"errors"
	"strings"

	//"github.com/badoux/checkmail"
	"gorm.io/gorm"
)

type Teacher struct {
	gorm.Model
	Username string `gorm:"size:50;unique_index;not null"    json:"username"`
	Email    string `gorm:"size:100;unique_index;not null"   json:"email"`
	Password string `gorm:"size:100;not null"                json:"password"`
	//Notification []Notification
}

// func BeforeSave hashes teacher password
func (t *Teacher) BeforeSave() error {
	password := strings.TrimSpace(t.Password)
	hashedpassword, err := HashPassword(password)
	if err != nil {
		return err
	}
	t.Password = string(hashedpassword)
	return nil
}

// func Prepare strips white spaces off input
func (t *Teacher) Prepare() {
	t.Email = strings.TrimSpace(t.Email)
	t.Password = strings.TrimSpace(t.Password)
	t.Username = strings.TrimSpace(t.Username)
}

// validates teacher input
func (t *Teacher) Validate(action string) error {
	switch strings.ToLower("action") {
	case "login":
		if t.Email == "" {
			return errors.New("please enter email")
		}
		if t.Password == "" {
			return errors.New("please enter password")
		}
		return nil

		// default: //all fields required
		// 	if t.Username == "" {
		// 		return errors.New("please enter username")
		// 	}
		// 	if t.Email == "" {
		// 		return errors.New("please enter email")
		// 	}
		// 	if t.Password == "" {
		// 		return errors.New("please enter password")
		// 	}
		// 	if err := checkmail.ValidateFormat(t.Email); err != nil {
		// 		return errors.New("invalid email")
		// 	}
		// 	return nil
	}
	return nil
}

// func Saveteacher adds teacher to database
func (t *Teacher) SaveTeacher(db *gorm.DB) (*Teacher, error) {
	//var err error
	err := db.Debug().Create(&t).Error
	if err != nil {
		return &Teacher{}, err
	}
	return t, nil
}

// func GetTeacher gets teacher from database
func (t *Teacher) GetTeacher(db *gorm.DB) (*Teacher, error) {
	account := &Teacher{}
	if err := db.Debug().Table("teachers").Where("email=?", t.Email).First(account).Error; err != nil {
		return nil, err
	}
	return account, nil
}

// func indAllTeachers gets all teachers
func (t *Teacher) FindAllTeachers(db *gorm.DB) (*[]Teacher, error) {
	accounts := &[]Teacher{}
	if err := db.Debug().Table("teachers").Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}
