package models

type Book struct {
	Subject   string `gorm:"size:50;not null" json:"subject"`
	StudentId uint
	IsRead    bool `gorm:"size:50"          json:"isread"`
}
