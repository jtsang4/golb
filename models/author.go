package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type BasicAuthor struct {
	Username string
	Email    string
	Password string
	Name     string
}

type Author struct {
	Id          int64
	CreatedTime time.Time
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

func AddOneAuthor(db *sql.DB, a BasicAuthor) (author Author, err error) {
	allAuthors, err := GetAllAuthors(db)
	if err != nil {
		return author, err
	}
	for _, author := range allAuthors {
		if author.Username == a.Username {
			return author, errors.New(fmt.Sprintf("Add author failed, the username %s is existing.", a.Username))
		} else if author.Email == a.Email {
			return author, errors.New(fmt.Sprintf("Add author failed, the email %s is existing.", a.Email))
		}
	}
	currentTime := time.Now()
	result, err := db.Exec(
		"INSERT INTO author (username, email, password, name, created_time) VALUES ($1, $2, $3, $4, $5)",
		a.Username, a.Email, a.Password, a.Name, currentTime,
	)
	if err != nil {
		return author, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return author, err
	}
	author = Author{
		Id:          id,
		CreatedTime: currentTime,
		BasicAuthor: a,
	}
	return author, nil
}

func GetAllAuthors(db *sql.DB) (authors []Author, err error) {
	rows, err := db.Query("SELECT id, username, email, password, name, created_time FROM author")
	if err != nil {
		return authors, err
	}
	for rows.Next() {
		var id int64
		var username string
		var email string
		var password string
		var name string
		var createdTime time.Time
		rows.Scan(&id, &username, &email, &password, &name, &createdTime)
		a := BasicAuthor{
			Username: username,
			Email:    email,
			Password: password,
			Name:     name,
		}
		authors = append(authors, Author{
			Id:          id,
			CreatedTime: createdTime,
			BasicAuthor: a,
		})
	}
	return authors, nil
}

func GetOneAuthor(db *sql.DB, condition ...string) (author Author, err error) {
	query := "SELECT id, username, email, password, name, created_time FROM author"
	if len(condition) > 0 {
		query += fmt.Sprintf(" %s", condition[0])
	}
	row := db.QueryRow(query)
	var id int64
	var username string
	var email string
	var password string
	var name string
	var createdTime time.Time
	err = row.Scan(&id, &username, &email, &password, &name, &createdTime)
	if err != nil {
		return author, err
	}
	basicAuthor := BasicAuthor{
		Username: username,
		Email:    email,
		Password: password,
		Name:     name,
	}
	author = Author{
		Id:          id,
		CreatedTime: createdTime,
		BasicAuthor: basicAuthor,
	}
	return author, nil
}

func GetAuthorById(db *sql.DB, id int64) (Author, error) {
	condition := fmt.Sprintf("WHERE id = %s", id)
	return GetOneAuthor(db, condition)
}

func GetAuthorByEmail(db *sql.DB, email string) (Author, error) {
	condition := fmt.Sprintf("WHERE email = %s", email)
	return GetOneAuthor(db, condition)
}

func UpdateOneAuthor(db *sql.DB, id int64, a BasicAuthor) (author Author, err error) {
	authors, err := GetAllAuthors(db)
	if err != nil {
		return author, err
	}
	for _, author := range authors {
		if author.Username == a.Username {
			return author, errors.New(fmt.Sprintf("Update author info failed, the username %s is duplicated.", a.Username))
		} else if author.Email == a.Email {
			return author, errors.New(fmt.Sprintf("Update author info failed, the email %s is duplicated.", a.Email))
		}
	}
	query := fmt.Sprintf(
		"UPDATE author SET username = %s, email = %s， password = %s, name = %s WHERE id = %d",
		a.Username, a.Email, a.Password, a.Name, id,
	)
	_, err = db.Exec(query)
	if err != nil {
		return author, err
	}
	return GetAuthorById(db, id)
}

func DeleteOneAuthor(db *sql.DB, id int64) (author Author, err error) {
	author, err = GetAuthorById(db, id)
	if err != nil {
		return author, err
	}
	query := fmt.Sprintf("DELETE FROM author WHERE id = %d", id)
	_, err = db.Exec(query)
	if err != nil {
		return author, err
	}
	return author, nil
}
