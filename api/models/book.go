package models

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Subject   string `gorm:"size:50;not null" json:"subject"`
	StudentId uint
	IsRead    bool `gorm:"size:50"          json:"isread"`
}

//Prepare strips off white spaces
func(b *Book) Prepare(){
	b.Subject = strings.TrimSpace(b.Subject)
	
}

//Validate Book input
func (b *Book) Validate(action string)error{
	switch strings.ToLower("action"){
	case "createbook":
		if b.Subject == ""{
			return errors.New("please input subject")
		}
		return nil

	case  "updatebook":
		if b.Subject == ""{
			return errors.New("please input subject")
		}
		return nil

	case "issuebook":
		if b.StudentId== 0{
			return errors.New("please input student id")
		}

	case "updatereadingprogress":
		if !b.IsRead{
			return errors.New("please update reading progress")
		}
	}
	return nil
	
}

//func CreateBook creates book
func (b *Book) CreateBook(db *gorm.DB)(*Book, error){
	//var err error
	err:= db.Debug().Create(&b).Error
	if err!= nil{
		return &Book{}, err
	}
	return b, nil
}

//func GetBook gets book from database
func (b *Book) GetBook(db *gorm.DB)(*Book, error){
	book:= &Book{}
	if err:= db.Debug().Table("books").Where("subject= ?", b.Subject).Or("booknumber= ?", b.ID).First(book).Error; err != nil{
		return nil, err
	}
	return book, nil
}
