package models

import (
	"database/sql"
	"fmt"
	"time"
)

type BasicPost struct {
	Title        string
	Content      string
	AuthorId     int64
	AuthorName   string
	CategoryId   int64
	CategoryName string
}

type Post struct {
	Id          int64
	CreatedTime time.Time
	UpdatedTime time.Time
	BasicPost
}

func CreatePostTable() (sql.Result, error) {
	return db.Exec(`
    CREATE TABLE IF NOT EXISTS post (
      id SERIAL PRIMARY KEY ,
      title text NOT NULL ,
      content text NULL ,
      author_id integer NOT NULL REFERENCES author (id) ,
      author_name text NOT NULL ,
      category_id integer NOT NULL REFERENCES category (id) ,
      category_name category NOT NULL,
      created_time timestamp NOT NULL ,
      updated_time timestamp NOT NULL
    )
  `)
}

func AddOnePost(p BasicPost) (post Post, err error) {
	currentTime := time.Now()
	var id int64
	err = db.QueryRow(
		"INSERT INTO post( title, content, author_id, author_name, category_id, category_name, created_time, updated_time ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id",
		p.Title, p.Content, p.CategoryId, p.CategoryName, p.AuthorId, p.AuthorName, currentTime, currentTime,
	).Scan(&id)
	if err != nil {
		return post, err
	}
	post = Post{
		Id:          id,
		CreatedTime: currentTime,
		UpdatedTime: currentTime,
		BasicPost:   p,
	}
	return post, nil
}

func GetPostsWithCondition(condition ...string) (posts []Post, err error) {
	query := "SELECT id, title, content, author_id, author_name, category_id, category_name, created_time, updated_time FROM post"
	if len(condition) > 0 {
		query += fmt.Sprintf(" %s", condition[0])
	}
	rows, err := db.Query(query)
	if err != nil {
		return posts, err
	}
	for rows.Next() {
		var (
			id           int64
			title        string
			content      string
			authorId     int64
			authorName   string
			categoryId   int64
			categoryName string
			createdTime  time.Time
			updatedTime  time.Time
		)
		err = rows.Scan(&id, &title, &content, &authorId, &authorName, &categoryId, &categoryName, &createdTime, &updatedTime)
		if err != nil {
			return posts, err
		}
		p := BasicPost{
			Title:        title,
			Content:      content,
			AuthorId:     authorId,
			AuthorName:   authorName,
			CategoryId:   categoryId,
			CategoryName: categoryName,
		}
		posts = append(posts, Post{
			Id:          id,
			CreatedTime: createdTime,
			UpdatedTime: updatedTime,
			BasicPost:   p,
		})
	}
	return posts, nil
}

func GetAllPosts() ([]Post, error) {
	return GetPostsWithCondition()
}

func GetPostsByAuthorId(authorId int64) (posts []Post, err error) {
	condition := fmt.Sprintf("WHERE author_id = %d", authorId)
	return GetPostsWithCondition(condition)
}

func GetPostsByCategoryId(categoryId int64) (posts []Post, err error) {
	condition := fmt.Sprintf("WHERE category_id = %d", categoryId)
	return GetPostsWithCondition(condition)
}

func GetOnePostWithCondition(condition ...string) (post Post, err error) {
	query := "SELECT id, title, content, author_id, author_name, category_id, category_name, created_time, updated_time FROM post"
	if len(condition) > 0 {
		query += fmt.Sprintf(" %s", condition[0])
	}
	row := db.QueryRow(query)
	var (
		id           int64
		title        string
		content      string
		authorId     int64
		authorName   string
		categoryId   int64
		categoryName string
		createdTime  time.Time
		updatedTime  time.Time
	)
	err = row.Scan(&id, &title, &content, &authorId, &authorName, &categoryId, &categoryName, &createdTime, &updatedTime)
	if err != nil {
		return post, err
	}
	p := BasicPost{
		Title:        title,
		Content:      content,
		AuthorId:     authorId,
		AuthorName:   authorName,
		CategoryId:   categoryId,
		CategoryName: categoryName,
	}
	post = Post{
		Id:          id,
		CreatedTime: createdTime,
		UpdatedTime: updatedTime,
		BasicPost:   p,
	}
	return post, nil
}

func GetOnePostById(id int64) (Post, error) {
	condition := fmt.Sprintf("WHERE id = %d", id)
	return GetOnePostWithCondition(condition)
}

func UpdateOnePost(id int64, p BasicPost) (post Post, err error) {
	var createdTime time.Time
	currentTime := time.Now()
	err = db.QueryRow(
		"UPDATE post SET title = $1, content = $2, category_id = $3, category_name = $4, updated_time = $5 WHERE id = $6 RETURNING created_time",
		p.Title, p.Content, p.CategoryId, p.CategoryName, currentTime, id,
	).Scan(&createdTime)
	if err != nil {
		return post, err
	}
	post = Post{
		Id:          id,
		CreatedTime: createdTime,
		UpdatedTime: currentTime,
		BasicPost:   p,
	}
	return post, nil
}

func DeleteOnePost(id int64) (post Post, err error) {
	p, err := GetOnePostById(id)
	if err != nil {
		return post, err
	}
	_, err = db.Exec("DELETE FROM post WHERE id = $1", id)
	if err != nil {
		return post, err
	}
	return p, nil
}
