package models

import (
	"database/sql"
	"log"
	"time"
)

type BasicPost struct {
	Title       string
	AuthorId    int64
	AuthorName  string
	Content     string
	CreatedTime time.Time
	UpdatedTime time.Time
}

type Post struct {
	Id int64
	BasicPost
}

func CreatePostTable(db *sql.DB) (sql.Result, error) {
	return db.Exec(`
    CREATE TABLE IF NOT EXISTS post (
      id SERIAL PRIMARY KEY ,
      title text NOT NULL ,
      content text NULL ,
      author_id integer NOT NULL REFERENCES author (id) ,
      author_name text NOT NULL ,
      category_id integer NOT NULL REFERENCES category (id) ,
      created_time timestamp NOT NULL ,
      updated_time timestamp NOT NULL
    )
  `)
}

func AddOnePost(db *sql.DB, p BasicPost) {
	_, err := db.Exec(
		"INSERT INTO post( title, author_id, author_name, content, created_time, updated_time ) VALUES ($1, $2, $3, $4, $5, $6)",
		p.Title, p.AuthorId, p.AuthorName, p.Content, p.CreatedTime, p.UpdatedTime,
	)
	if err != nil {
		log.Fatalln(err)
	}
}
