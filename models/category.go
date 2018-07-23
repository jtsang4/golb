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

func CreateCategoryTable() (sql.Result, error) {
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

func AddOneCategory(c BasicCategory) (category Category, err error) {
	categories, err := GetAllCategories()
	if err != nil {
		return category, err
	}
	for _, cat := range categories {
		if cat.AuthorId == c.AuthorId && cat.Title == c.Title {
			return category, errors.New(fmt.Sprintf("Add category failed, the title %s is existing.", c.Title))
		}
	}
	currentTime := time.Now()
	var id int64
	err = db.QueryRow(
		"INSERT INTO category( title, author_id, author_name, created_time, updated_time ) VALUES ('$1', $2, '$3', $4, $5) RETURNING id",
		c.Title, c.AuthorId, c.AuthorName, currentTime, currentTime,
	).Scan(&id)
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

func GetCategoriesWithCondition(condition ...string) (categories []Category, err error) {
	query := "SELECT id, title, author_id, author_name, created_time, updated_time FROM category"
	if len(condition) > 0 {
		query += fmt.Sprintf(" %s", condition[0])
	}
	rows, err := db.Query(query)
	if err != nil {
		return categories, err
	}
	for rows.Next() {
		var (
			id          int64
			title       string
			authorId    int64
			authorName  string
			createdTime time.Time
			updatedTime time.Time
		)
		err = rows.Scan(&id, &title, &authorId, &authorName, &createdTime, &updatedTime)
		if err != nil {
			return categories, err
		}
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

func GetAllCategories() (categories []Category, err error) {
	return GetCategoriesWithCondition()
}

func GetCategoriesByAuthorId(authorId int64) (categories []Category, err error) {
	condition := fmt.Sprintf(" WHERE author_id = %d", authorId)
	return GetCategoriesWithCondition(condition)
}

func GetOneCategoryWithCondition(condition ...string) (category Category, err error) {
	query := "SELECT id, title, author_id, author_name, created_time, updated_time FROM category"
	if len(condition) > 0 {
		query += fmt.Sprintf(" %s", condition[0])
	}
	row := db.QueryRow(query)
	var (
		id          int64
		title       string
		authorId    int64
		authorName  string
		createdTime time.Time
		updatedTime time.Time
	)
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

func GetOneCategoryById(id int64) (category Category, err error) {
	condition := fmt.Sprintf("WHERE id = %d", id)
	return GetOneCategoryWithCondition(condition)
}

func GetOneCategoryByTitle(title string) (category Category, err error) {
	condition := fmt.Sprintf("WHERE title = '%s'", title)
	return GetOneCategoryWithCondition(condition)
}

func UpdateOneCategory(id int64, title string) (category Category, err error) {
	categories, err := GetCategoriesByAuthorId(id)
	if err != nil {
		return category, err
	}
	for _, cat := range categories {
		if cat.AuthorId == id && cat.Title == title {
			return category, errors.New(fmt.Sprintf("Update category info failed, the title %s is existing", title))
		}
	}
	var (
		authorId    int64
		authorName  string
		createdTime time.Time
	)
	currentTime := time.Now()
	err = db.QueryRow(
		"UPDATE category SET title = '$1', updated_time = $2 WHERE id = $3 RETURNING author_id, author_name, created_time",
		title, currentTime, id,
	).Scan(&authorId, &authorName, &createdTime)
	if err != nil {
		return category, err
	}
	category = Category{
		Id:          id,
		CreatedTime: createdTime,
		UpdatedTime: currentTime,
		BasicCategory: BasicCategory{
			Title:      title,
			AuthorId:   authorId,
			AuthorName: authorName,
		},
	}
	return category, nil
}

func DeleteOneCategory(id int64) (category Category, err error) {
	c, err := GetOneCategoryById(id)
	if err != nil {
		return category, err
	}
	_, err = db.Exec("DELETE FROM category WHERE id = $1", id)
	if err != nil {
		return category, err
	}
	return c, nil
}
