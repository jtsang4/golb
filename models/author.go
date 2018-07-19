package models

import (
	"database/sql"
	"log"
	"time"
)

type BasicAuthor struct {
	Username    string
	Email       string
	Password    string
	Name        string
	CreatedTime time.Time
}

type Author struct {
	Id uint32
	BasicAuthor
}

func CreateAuthorTable(db *sql.DB) (sql.Result, error) {
	return db.Exec(`
    CREATE TABLE IF NOT EXISTS author (
      id SERIAL PRIMARY KEY ,
      username text UNIQUE NOT NULL ,
      email text UNIQUE NOT NULL ,
      password text NOT NULL ,
      name text NOT NULL ,
      created_time timestamp NOT NULL
    )
  `)
}

func AddOneAuthor(db *sql.DB, a BasicAuthor) {
	allAuthors := GetAllAuthors(db)
	for _, author := range allAuthors {
		if author.Username == a.Username || author.Email == a.Email {
			return
		}
	}
	_, err := db.Exec(
		"INSERT INTO author (username, email, password, name, created_time) VALUES ($1, $2, $3, $4, $5)",
		a.Username, a.Email, a.Password, a.Name, a.CreatedTime,
	)
	if err != nil {
		log.Fatalln(err)
	}
}

func GetOneAuthor(db *sql.DB) Author {
	row := db.QueryRow("SELECT id, username, email, password, name, created_time FROM author")
	var id uint32
	var username string
	var email string
	var password string
	var name string
	var createdTime time.Time
	row.Scan(&id, &username, &email, &password, &name, &createdTime)
	basicAuthor := BasicAuthor{
		Username:    username,
		Email:       email,
		Password:    password,
		Name:        name,
		CreatedTime: createdTime,
	}
	return Author{
		Id:          id,
		BasicAuthor: basicAuthor,
	}
}

func GetAllAuthors(db *sql.DB) (authors []Author) {
	rows, err := db.Query("SELECT id, username, email, password, name, created_time FROM author")
	if err != nil {
		log.Fatalln(err)
	}
	for rows.Next() {
		var id uint32
		var username string
		var email string
		var password string
		var name string
		var createdTime time.Time
		rows.Scan(&id, &username, &email, &password, &name, &createdTime)
		a := BasicAuthor{
			Username:    username,
			Email:       email,
			Password:    password,
			Name:        name,
			CreatedTime: createdTime,
		}
		authors = append(authors, Author{
			Id:          id,
			BasicAuthor: a,
		})
	}
	return authors
}
