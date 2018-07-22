package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type BasicCategory struct {
	Title      string
	AuthorId   int64
	AuthorName string
}

type Category struct {
	Id          int64
	CreatedTime time.Time
	UpdatedTime time.Time
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

func AddOneCategory(db *sql.DB, c BasicCategory) (category Category, err error) {
	categories, err := GetAllCategories(db)
	if err != nil {
		return category, err
	}
	for _, cat := range categories {
		if cat.AuthorId == c.AuthorId && cat.Title == c.Title {
			return category, errors.New(fmt.Sprintf("Add category failed, the title %s is existing.", c.Title))
		}
	}
	currentTime := time.Now()
	result, err := db.Exec(
		"INSERT INTO category( title, author_id, author_name, created_time, updated_time ) VALUES ($1, $2, $3, $4, $5)",
		c.Title, c.AuthorId, c.AuthorName, currentTime, currentTime,
	)
	if err != nil {
		return category, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return category, err
	}
	category = Category{
		Id:            id,
		CreatedTime:   currentTime,
		UpdatedTime:   currentTime,
		BasicCategory: c,
	}
	return category, nil
}

func GetCategoriesWithCondition(db *sql.DB, condition ...string) (categories []Category, err error) {
	query := "SELECT id, title, author_id, author_name, created_time, updated_time FROM category"
	if len(condition) > 0 {
		query += fmt.Sprintf(" %s", condition)
	}
	rows, err := db.Query(query)
	if err != nil {
		return categories, err
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
			Title:      title,
			AuthorId:   authorId,
			AuthorName: authorName,
		}
		categories = append(categories, Category{
			Id:            id,
			CreatedTime:   createdTime,
			UpdatedTime:   updatedTime,
			BasicCategory: c,
		})
	}
	return categories, nil
}

func GetAllCategories(db *sql.DB) (categories []Category, err error) {
	return GetCategoriesWithCondition(db)
}

func GetOneCategory(db *sql.DB, condition ...string) (category Category, err error) {
	query := "SELECT id, title, author_id, author_name, created_time, updated_time FROM category"
	if len(condition) > 0 {
		query += fmt.Sprintf(" %s", condition[0])
	}
	row := db.QueryRow(query)
	var id int64
	var title string
	var authorId int64
	var authorName string
	var createdTime time.Time
	var updatedTime time.Time
	err = row.Scan(&id, &title, &authorId, &authorName, &createdTime, &updatedTime)
	if err != nil {
		return category, err
	}
	basicCategory := BasicCategory{
		Title:      title,
		AuthorId:   authorId,
		AuthorName: authorName,
	}
	category = Category{
		Id:            id,
		CreatedTime:   createdTime,
		UpdatedTime:   updatedTime,
		BasicCategory: basicCategory,
	}
	return category, nil
}

func GetOneCategoryById(db *sql.DB, id int64) (category Category, err error) {
	condition := fmt.Sprintf("WHERE id = %d", id)
	return GetOneCategory(db, condition)
}

func GetOneCategoryByTitle(db *sql.DB, title string) (category Category, err error) {
	condition := fmt.Sprintf("WHERE title = %s", title)
	return GetOneCategory(db, condition)
}

func GetCategoriesByAuthorId(db *sql.DB, authorId int64) (categories []Category, err error) {
	condition := fmt.Sprintf(" WHERE author_id = %d", authorId)
	return GetCategoriesWithCondition(db, condition)
}

func UpdateOneCategory(db *sql.DB, c Category) (category Category, err error) {
	categories, err := GetCategoriesByAuthorId(db, c.AuthorId)
	if err != nil {
		return category, err
	}
	for _, cat := range categories {
		if cat.AuthorId == c.AuthorId && c.Title == cat.Title {
			return category, errors.New(fmt.Sprintf("Update category info failed, the title %s is existing", c.Title))
		}
	}
	query := fmt.Sprintf(
		"UPDATE category SET title = %s WHERE id = %d",
		c.Title, c.Id,
	)
	_, err = db.Exec(query)
	if err != nil {
		return category, err
	}
	return c, nil
}

func DeleteOneCategory(db *sql.DB, id int64) (category Category, err error) {
	category, err = GetOneCategoryById(db, id)
	if err != nil {
		return category, err
	}
	query := fmt.Sprintf("DELETE FROM category WHERE id = %d", id)
	_, err = db.Exec(query)
	if err != nil {
		return category, err
	}
	return category, nil
}
