package models

import (
	"database/sql"
	"fmt"
	"log"
)

func InsertOneCategory(db *sql.DB) {
	author := BasicAuthor{
		Username: "wtzeng",
		Email:    "wtzeng@ss.com",
		Password: "gogoooo11",
		Name:     "wentao",
	}

	AddOneAuthor(db, author)
	at, err := GetOneAuthor(db)
	if err != nil {
		log.Fatal(err)
	}
	category := BasicCategory{
		Title:      "the first category",
		AuthorId:   at.Id,
		AuthorName: at.Name,
	}
	AddOneCategory(db, category)
	fmt.Println(GetAllCategories(db))
}

func TestStart(db *sql.DB) {
	InsertOneCategory(db)
}
