package models

import (
	"database/sql"
	"log"
	"time"
)

type BasicCategory struct {
	Title       string
	AuthorId    int64
	AuthorName  string
	CreatedTime time.Time
	UpdatedTime time.Time
}

type Category struct {
	Id int64
	BasicCategory
}

func CreateCategoryTable(db *sql.DB) (sql.Result, error) {
	return db.Exec(`
    CREATE TABLE IF NOT EXISTS category (
      id SERIAL PRIMARY KEY ,
      title text NOT NULL ,
      author_id integer NOT NULL REFERENCES author (id) ,
      author_name text NOT NULL ,
      created_time timestamp NOT NULL ,
      updated_time timestamp NOT NULL ,
      UNIQUE (title, author_id)
    )
  `)
}

func AddOneCategory(db *sql.DB, c BasicCategory) {
	_, err := db.Exec(
		"INSERT INTO category( title, author_id, author_name, created_time, updated_time ) VALUES ($1, $2, $3, $4, $5)",
		c.Title, c.AuthorId, c.AuthorName, c.CreatedTime, c.UpdatedTime,
	)
	if err != nil {
		log.Fatalln(err)
	}
}

func GetAllCategories(db *sql.DB) (categories []Category) {
	rows, err := db.Query("SELECT id, title, author_id, author_name, created_time, updated_time FROM category")
	if err != nil {
		log.Fatalln(err)
	}
	for rows.Next() {
		var id int64
		var title string
		var authorId int64
		var authorName string
		var createdTime time.Time
		var updatedTime time.Time
		rows.Scan(&id, &title, &authorId, &authorName, &createdTime, &updatedTime)
		c := BasicCategory{
			Title:       title,
			AuthorId:    authorId,
			AuthorName:  authorName,
			CreatedTime: createdTime,
			UpdatedTime: updatedTime,
		}
		categories = append(categories, Category{
			Id:            id,
			BasicCategory: c,
		})
	}
	return categories
}
