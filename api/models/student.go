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
	Username    string `gorm:"size:100;not null"              json:"username"`
	PhoneNumber int    `gorm:"size:20;not null"               json:"phonenumber"`
	Email       string `gorm:"type:varchar(100);unique_index" json:"email"`
	Password    string `gorm:"size:100;not null"              json:"password"`
	Books       int    `json:"books" `
	//Notification Notification

}

type StudentAndBooks struct {
	StudentUsername string `gorm:"references:Username"`
	BookCount       uint
}

type book_count struct {
	bookCount uint
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
	//s.ProfileImage = strings.TrimSpace(s.ProfileImage)
}

// Validate user input
func (s *Student) Validate(action string) error {
	switch strings.ToLower(action) {
	case "login":
		if s.Email == "" {
			return errors.New("please enter email address")
		}
		if s.Password == "" {
			return errors.New("please provide password to login")
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
	//var err error
	err := db.Debug().Create(&s).Error
	if err != nil {
		return &Student{}, err
	}
	return s, nil

}

// Get  Student based on email or phone number
func (s *Student) GetStudent(db *gorm.DB) (*Student, error) {
	account := &Student{}
	if err := db.Debug().Table("students").Where("email= ?", s.Email).First(account).Error; err != nil {
		return nil, err
	}
	return account, nil
}

// function GetStudentById gets specific student by their id
func (s *Student) GetStudentById(id int, db *gorm.DB) (*Student, error) {
	user := &Student{}

	if err := db.Debug().Table("students").Where("id= ?", id).First(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// function UpdateUser updates specific user
func (s *Student) UpdateStudent(id int, db *gorm.DB) (*Student, error) {
	if err := db.Debug().Table("students").Where("id= ?", s.ID).Updates(Student{
		Username:    s.Username,
		Email:       s.Email,
		Password:    s.Password,
		PhoneNumber: s.PhoneNumber}).Error; err != nil {
		return &Student{}, err
	}
	return s, nil

}

// func GetStudents gets all students in db
func GetStudents(db *gorm.DB) (*[]Student, error) {
	users := []Student{}
	if err := db.Debug().Table("students").Find(&users).Error; err != nil {
		return &[]Student{}, err
	}
	return &users, nil
}

// func GetStudentsAndBooks gets students with assigned books
func GetStudentsAndBooks(db *gorm.DB) (*[]Student, error) {
	users := []Student{}
	if err := db.Debug().Preload("Books").
		Joins("INNER JOIN books on books.student_id = students.id").Find(&users).Error; err != nil {
		return &[]Student{}, err
	}
	return &users, nil
}

// func CountBooks counts books group by username
func (s *Student) CountBooks(db *gorm.DB) (*[]StudentAndBooks, error) {
	students := []StudentAndBooks{}
	//var count int64
	if err := db.Debug().Raw("select students.username as student_username ,count(books.subject) as book_count from students inner join books on books.student_id= students.id group by students.id").Scan(&students).Error; err != nil {
		return &[]StudentAndBooks{}, err
	}
	return &students, nil
}

// func PopulateBooks adds books assigned to student
func (s *Student) PopulateBooks(studentID int, db *gorm.DB) (*Student, error) {
	//books := []Book{}
	students := Student{}
	//var count int64
	if err := db.Debug().Raw("select students.id,students.created_at,students.updated_at,students.deleted_at,students.username,students.phone_number,students.email,students.password, count(*) as books FROM students inner join books on books.student_id= students.id WHERE students.id= ? AND students.deleted_at IS NULL group by students.id", studentID).Scan(&students).Error; err != nil {
		return nil, err
	}
	return &students, nil
}


