package models

type Book struct {
	ID     uint   `gorm:"primaryKey"`
	Title  string `gorm:"unique;not null"`
	Author string `gorm:"not null"`
}
