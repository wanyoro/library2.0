package models

import (
	"errors"
	"strings"

	"gorm.io/gorm"
	//"fmt"
	//"lib2.0/api/main"
)

type Book struct {
	gorm.Model
	Subject string `gorm:"size:50;not null" json:"subject"`
	//
	IsRead    *bool `gorm:"size:50"          json:"isread"`
	StudentID uint  //`gorm:"foreignKey:StudentID"`
	TeacherID uint
}

// Prepare strips off white spaces
func (b *Book) Prepare() {
	b.Subject = strings.TrimSpace(b.Subject)

}

// Validate Book input
func (b *Book) Validate(action string) error {
	switch strings.ToLower("action") {
	case "createbook":
		if b.Subject == "" {
			return errors.New("please input subject")
		}
		return nil

	case "updatebook":
		if b.Subject == "" {
			return errors.New("please input subject")
		}
		return nil

		// case "issuebook":
		// 	if b.StudentID == 0 {
		// 		return errors.New("please input student id")
		// 	}

		// case "updatereadingprogress":
		// 	if !b.IsRead {
		// 		return errors.New("please update reading progress")
		// 	}
	}
	return nil

}

// func CreateBook creates book
func (b *Book) CreateBook(db *gorm.DB) (*Book, error) {
	//var err error
	err := db.Debug().Create(&b).Error
	if err != nil {
		return &Book{}, err
	}
	return b, nil
}

//func

// func GetBook gets book from database
func (b *Book) GetBookById(id int, db *gorm.DB) (*Book, error) {
	book := &Book{}
	if err := db.Debug().Table("books").Where("id= ?", id).First(book).Error; err != nil {
		return nil, err
	}
	return book, nil
}

// func GetBooks gets all books from database
func (b *Book) GetBooks(db *gorm.DB) (*[]Book, error) {
	books := []Book{}
	if err := db.Debug().Table("books").Find(&books).Error; err != nil {
		return nil, err
	}
	return &books, nil
}

// func UpdateBook updates book subject and student assigned
func (b *Book) UpdateBook(id int, db *gorm.DB) (*Book, error) {
	//var resp = map[string]interface{}{"status": "successful", "message": "book updated successfully"}
	BookSubject := b.Subject
	BookId := b.ID
	notif := Notification{}
	if err := db.Debug().Table("books").Where("id =?", b.ID).Updates(Book{
		//Subject: b.Subject,

		StudentID: b.StudentID,
		IsRead:    b.IsRead,
		TeacherID: b.TeacherID}).Error; err != nil {
		return &Book{}, err
	}
	if *b.IsRead {
		if db.Debug().Table("notifications").Where("book_subject= ?", BookSubject).First(&notif).RowsAffected == 0 {
			db.Create(&Notification{
				Message:     "Book is read",
				StudentID:   b.StudentID,
				BookSubject: BookSubject,
				BookId:      BookId,
				TeacherID:   b.TeacherID,
			})
		}
	}

	return b, nil
}

// func UpdateReadingProgress sets reading progress to true and creates notif
func (b *Book) UpdateReadingProgress(id int, db *gorm.DB) (*Book, error) {
	if err := db.Debug().Table("books").Where("id= ?", b.ID).Updates(Book{
		IsRead: b.IsRead}).Error; err != nil {
		return &Book{}, err
	}
	return b, nil
}

//func checks reading progress and returns boolean
// func (b *Book) CheckReadingProgress (id int, db *gorm.DB) bool {
// 	if err :=
// }
