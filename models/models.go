package models

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

var db *sql.DB

func init() {
	registerDB()
}

func registerDB() {
	var err error
	db, err = sql.Open("postgres", "user=postgres dbname=golb password= sslmode=disable")
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	err = createTables()
	if err != nil {
		log.Fatalln(err)
	}
}

func CloseDB() {
	db.Close()
}

func createTables() error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = CreateAuthorTable()
	if err != nil {
		return err
	}
	_, err = CreateCategoryTable()
	if err != nil {
		return err
	}
	_, err = CreatePostTable()
	if err != nil {
		return err
	}

	return tx.Commit()
}
