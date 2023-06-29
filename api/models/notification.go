package models

type Notification struct {
	Message   string `gorm:"size:200;not null" json:"message"`
	StudentID uint
	TeacherID uint
}
