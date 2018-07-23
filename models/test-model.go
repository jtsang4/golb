package models

import (
	"fmt"
	"log"
)

func InsertOneCategory() {
	author := BasicAuthor{
		Username: "wtzeng",
		Email:    "wtzeng@ss.com",
		Password: "gogoooo11",
		Name:     "wentao",
	}

	AddOneAuthor(author)
	at, err := GetOneAuthorWithCondition()
	if err != nil {
		log.Fatal(err)
	}
	category := BasicCategory{
		Title:      "the first category",
		AuthorId:   at.Id,
		AuthorName: at.Name,
	}
	AddOneCategory(category)
	fmt.Println(GetAllCategories())
}

func TestStart() {
	InsertOneCategory()
}
