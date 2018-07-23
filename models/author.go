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

func CreateAuthorTable() (sql.Result, error) {
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

func AddOneAuthor(a BasicAuthor) (author Author, err error) {
	allAuthors, err := GetAllAuthors()
	if err != nil {
		return author, err
	}
	for _, auth := range allAuthors {
		if auth.Username == a.Username {
			return author, errors.New(fmt.Sprintf("Add author failed, the username %s is existing.", a.Username))
		} else if auth.Email == a.Email {
			return author, errors.New(fmt.Sprintf("Add author failed, the email %s is existing.", a.Email))
		}
	}
	currentTime := time.Now()
	var id int64
	err = db.QueryRow(
		"INSERT INTO author (username, email, password, name, created_time) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		a.Username, a.Email, a.Password, a.Name, currentTime,
	).Scan(&id)
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

func GetAllAuthors() (authors []Author, err error) {
	rows, err := db.Query("SELECT id, username, email, password, name, created_time FROM author")
	if err != nil {
		return authors, err
	}
	for rows.Next() {
		var (
			id          int64
			username    string
			email       string
			password    string
			name        string
			createdTime time.Time
		)
		err = rows.Scan(&id, &username, &email, &password, &name, &createdTime)
		if err != nil {
			return authors, err
		}
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

func GetOneAuthorWithCondition(condition ...string) (author Author, err error) {
	query := "SELECT id, username, email, password, name, created_time FROM author"
	if len(condition) > 0 {
		query += fmt.Sprintf(" %s", condition[0])
	}
	row := db.QueryRow(query)
	var (
		id          int64
		username    string
		email       string
		password    string
		name        string
		createdTime time.Time
	)
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

func GetAuthorById(id int64) (Author, error) {
	condition := fmt.Sprintf("WHERE id = %d", id)
	return GetOneAuthorWithCondition(condition)
}

func GetAuthorByEmail(email string) (Author, error) {
	condition := fmt.Sprintf("WHERE email = '%s'", email)
	return GetOneAuthorWithCondition(condition)
}

func UpdateOneAuthor(id int64, a BasicAuthor) (author Author, err error) {
	authors, err := GetAllAuthors()
	if err != nil {
		return author, err
	}
	for _, author := range authors {
		if author.Username == a.Username {
			return author, errors.New(fmt.Sprintf("Update author info failed, the username %s is existing.", a.Username))
		} else if author.Email == a.Email {
			return author, errors.New(fmt.Sprintf("Update author info failed, the email %s is existing.", a.Email))
		}
	}
	var createdTime time.Time
	err = db.QueryRow(
		"UPDATE author SET username = $1, email = $2, password = $3, name = $4 WHERE id = $5 RETURNING created_time",
		a.Username, a.Email, a.Password, a.Name, id,
	).Scan(&createdTime)
	if err != nil {
		return author, err
	}
	author = Author{
		Id:          id,
		CreatedTime: createdTime,
		BasicAuthor: a,
	}
	return author, nil
}

func DeleteOneAuthor(id int64) (author Author, err error) {
	a, err := GetAuthorById(id)
	if err != nil {
		return author, err
	}
	_, err = db.Exec("DELETE FROM author WHERE id = $1", id)
	if err != nil {
		return author, err
	}
	return a, nil
}
