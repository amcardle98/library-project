package models

import "gorm.io/gorm"

type Book struct {
	ID    uint    `gorm:"primary key;autoIncrement" json:"id"`
	Title *string `json:"title"`
}

func Migration(db *gorm.DB) error {
	err := db.AutoMigrate(&Book{})
	return err
}
