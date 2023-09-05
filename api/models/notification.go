package models

import (
	"gorm.io/gorm"
)

type Notification struct {
	gorm.Model
	Message     string `gorm:"size:200;not null" json:"message"`
	StudentID   uint
	TeacherID   uint
	BookSubject string `gorm:"references:booksubject"`
	BookId      uint   `gorm:"references:bookID"`
	//TeacherUsername string `gorm:"foreignKey:teacherusername"`
}

type booksubject string

// func CreateNotif creates notification when reading progress is set to true
func (n *Notification) CreateNotif(db *gorm.DB) (*Notification, error) {
	err := db.Debug().Create(&n).Error
	if err != nil {
		return &Notification{}, err
	}
	return n, nil
}

// func GetNotifs gets all notifications from table
func (n *Notification) GetNotifs(db *gorm.DB) (*[]Notification, error) {
	notifs := []Notification{}
	if err := db.Debug().Table("notifications").Find(&notifs).Error; err != nil {
		return nil, err
	}
	return &notifs, nil
}

// func GetReadBooks gets books read by student
func (n *Notification) GetReadBooks(student_id int, db *gorm.DB) (*[]Notification, error) {
	//type book_subject string
	notif := &[]Notification{}
	if err := db.Debug().Table("notifications").Select("book_subject").Where("student_id=?", student_id).Find(notif).Error; err != nil {
		return nil, err
	}
	return notif, nil
}

//func GetNotifById
// func (n *Notification) GetNotifById (id int, db *gorm.DB)(*Notification, error){
// 	notif := Notification{}
// 	if err:= db.Debug().
// }
