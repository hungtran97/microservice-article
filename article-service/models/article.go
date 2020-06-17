package models

import (
	"log"
	"time"
)

type Article struct {
	ID          uint      `gorm:"primary_key"`
	Slug        string    `gorm:"column:slug;unique_index"`
	Title       string    `gorm:"column:title"`
	Description string    `gorm:"column:description;size:2048"`
	Body        string    `gorm:"column:body;size:2048"`
	AuthorID    uint      `gorm:"column:author_id"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (Article) TableName() string {
	return "articles"
}

func (db *DB) GetListArticle() ([]Article, error) {
	var articles []Article
	err := db.Table("articles").Find(&articles).Error

	if err != nil {
		log.Println("Errors: ", err)
	} else {
		log.Println("Articles: ", articles)
	}

	return articles, err
}
