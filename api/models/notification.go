package models

import (
	"gorm.io/gorm"
)

type Notification struct {
	Message     string `gorm:"size:200;not null" json:"message"`
	StudentID   uint
	TeacherID   uint
	BookSubject string `gorm:"references:booksubject"`
	//TeacherUsername string `gorm:"foreignKey:teacherusername"`
}

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
