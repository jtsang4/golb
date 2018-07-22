package models

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

var DB *sql.DB

func init() {
	registerDB()
}

func registerDB() {
	var err error
	DB, err = sql.Open("postgres", "user=postgres dbname=golb password= sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer DB.Close()
	err = DB.Ping()
	if err != nil {
		panic(err)
	}
	err = createTables(DB)
	if err != nil {
		log.Fatalln(err)
	}
	TestStart(DB)
}

func createTables(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = CreateAuthorTable(db)
	if err != nil {
		return err
	}
	_, err = CreateCategoryTable(db)
	if err != nil {
		return err
	}
	_, err = CreatePostTable(db)
	if err != nil {
		return err
	}

	return tx.Commit()
}
