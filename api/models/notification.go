package models

import(
	"gorm.io/gorm"
)

type Notification struct {
	Message   string `gorm:"size:200;not null" json:"message"`
	StudentID uint
	TeacherID uint
}

//func CreateNotif creates notification when reading progress is set to true
func (n *Notification) CreateNotif (db *gorm.DB) (*Notification, error){
	err := db.Debug().Create(&n).Error
	if err!= nil{
		return &Notification{}, err
	}
	return n, nil
}
