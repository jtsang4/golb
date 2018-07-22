package models

import (
	"database/sql"
	"fmt"
	"time"
)

type BasicPost struct {
	Title      string
	AuthorId   int64
	AuthorName string
	Content    string
}

type Post struct {
	Id          int64
	CreatedTime time.Time
	UpdatedTime time.Time
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

func AddOnePost(db *sql.DB, p BasicPost) (post Post, err error) {
	currentTime := time.Now()
	result, err := db.Exec(
		"INSERT INTO post( title, author_id, author_name, content, created_time, updated_time ) VALUES ($1, $2, $3, $4, $5, $6)",
		p.Title, p.AuthorId, p.AuthorName, p.Content, currentTime, currentTime,
	)
	if err != nil {
		return post, err
	}
	id, err := result.LastInsertId()
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

func GetPostsWithCondition(db *sql.DB, condition ...string) (posts []Post, err error) {
	query := "SELECT id, title, author_id, author_name, content, created_time, updated_time FROM post"
	if len(condition) > 0 {
		query += fmt.Sprintf(" %s", condition[0])
	}
	rows, err := db.Query(query)
	if err != nil {
		return posts, err
	}
	for rows.Next() {
		var id int64
		var title string
		var authorId int64
		var authorName string
		var content string
		var createdTime time.Time
		var updatedTime time.Time
		err = rows.Scan(&id, &title, &authorId, &authorName, &content, &createdTime, &updatedTime)
		if err != nil {
			return posts, err
		}
		p := BasicPost{
			Title:      title,
			AuthorId:   authorId,
			AuthorName: authorName,
			Content:    content,
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

func GetAllPosts(db *sql.DB) ([]Post, error) {
	return GetPostsWithCondition(db)
}

func GetPostsByAuthorId(db *sql.DB, authorId int64) (posts []Post, err error) {
	condition := fmt.Sprintf("WHERE author_id = %d", authorId)
	return GetPostsWithCondition(db, condition)
}

func GetOnePostWithCondition(db *sql.DB, condition ...string) (post Post, err error) {
	query := "SELECT id, title, author_id, author_name, content, created_time, updated_time FROM post"
	if len(condition) > 0 {
		query += fmt.Sprintf(" %s", condition[0])
	}
	row := db.QueryRow(query)
	var id int64
	var title string
	var authorId int64
	var authorName string
	var content string
	var createdTime time.Time
	var updatedTime time.Time
	err = row.Scan(&id, &title, &authorId, &authorName, &content, &createdTime, &updatedTime)
	if err != nil {
		return post, err
	}
	p := BasicPost{
		Title:      title,
		AuthorId:   authorId,
		AuthorName: authorName,
		Content:    content,
	}
	post = Post{
		Id:          id,
		CreatedTime: createdTime,
		UpdatedTime: updatedTime,
		BasicPost:   p,
	}
	return post, nil
}

func GetOnePostById(db *sql.DB, id int64) (Post, error) {
	condition := fmt.Sprintf("WHERE id = %d", id)
	return GetOnePostWithCondition(db, condition)
}

func UpdateOnePost(db *sql.DB, p Post) (post Post, err error) {
	currentTime := time.Now()
	query := fmt.Sprintf(
		"UPDATE post SET title = %s, content = %s, updated_time = %v WHERE id = %d",
		p.Title, p.Content, currentTime, p.Id,
	)
	_, err = db.Exec(query)
	if err != nil {
		return post, err
	}
	p.UpdatedTime = currentTime
	return p, nil
}

func DeleteOnePost(db *sql.DB, id int64) (post Post, err error) {
	p, err := GetOnePostById(db, id)
	if err != nil {
		return post, err
	}
	query := fmt.Sprintf("DELETE FROM post WHERE id = %d", id)
	_, err = db.Exec(query)
	if err != nil {
		return post, err
	}
	return p, nil
}
