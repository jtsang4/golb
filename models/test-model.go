package models

import (
	"database/sql"
	"fmt"
	"time"
)

func InsertOneCategory(db *sql.DB) {
	author := BasicAuthor{
		Username: "wtzeng",
		Email:    "wtzeng@ss.com",
		Password: "gogoooo11",
		Name:     "wentao",
	}

	AddOneAuthor(db, author)
	at := GetOneAuthor(db)
	category := BasicCategory{
		Title:       "the first category",
		AuthorId:    at.Id,
		AuthorName:  at.Name,
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}
	AddOneCategory(db, category)
	fmt.Println(GetAllCategories(db))
}

func TestStart(db *sql.DB) {
	InsertOneCategory(db)
}
