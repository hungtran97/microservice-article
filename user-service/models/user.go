package models

import (
	"log"
)

type User struct {
	ID           uint    `gorm:"primary_key"`
	Username     string  `gorm:"column:username"`
	Email        string  `gorm:"column:email;unique_index"`
	Bio          string  `gorm:"column:bio;size:1024"`
	Image        *string `gorm:"column:image"`
	PasswordHash string  `gorm:"column:password;not null"`
}

func (User) TableName() string {
	return "user_models"
}

func (db *DB) ListUser() ([]User, error) {
	var users []User
	err := db.Table("user_models").Find(&users).Error

	if err != nil {
		log.Println("Error: ", err)
	} else {
		log.Println("Users: ", users)
	}

	return users, nil
}
