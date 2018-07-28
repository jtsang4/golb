package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/wtzeng1/golb/models"
	"net/http"
	"strconv"
)

func RegisterPostController() {
	postRouter := router.Group("/post")

	{
		postRouter.POST("/", func(context *gin.Context) {
			title := context.PostForm("title")
			content := context.PostForm("content")
			authorId := context.PostForm("authorId")
			categoryId := context.PostForm("categoryId")
			atId, err := strconv.ParseInt(authorId, 10, 64)
			if err != nil {
				panic(err)
			}
			ctId, err := strconv.ParseInt(categoryId, 10, 64)
			if err != nil {
				panic(err)
			}
			author, err := models.GetAuthorById(atId)
			if err != nil {
				panic(err)
			}
			category, err := models.GetOneCategoryById(ctId)
			if err != nil {
				panic(err)
			}
			basicPost := models.BasicPost{
				Title:        title,
				Content:      content,
				AuthorId:     atId,
				AuthorName:   author.Name,
				CategoryId:   ctId,
				CategoryName: category.Title,
			}
			post, err := models.AddOnePost(basicPost)
			if err != nil {
				panic(err)
			}
			context.JSON(http.StatusOK, post)
		})

		postRouter.GET("/", func(context *gin.Context) {
			id := context.PostForm("id")
			authorId := context.PostForm("authorId")
			categoryId := context.PostForm("categoryId")
			if id == "" && authorId == "" && categoryId == "" {
				posts, err := models.GetAllPosts()
				if err != nil {
					panic(err)
				}
				context.JSON(http.StatusOK, posts)
			} else {
				if id != "" {
					postId, err := strconv.ParseInt(id, 10, 64)
					if err != nil {
						panic(err)
					}
					post, err := models.GetOnePostById(postId)
					if err != nil {
						panic(err)
					}
					context.JSON(http.StatusOK, post)
				} else {
					if authorId != "" && categoryId != "" {
						atId, err := strconv.ParseInt(authorId, 10, 64)
						if err != nil {
							panic(err)
						}
						ctId, err := strconv.ParseInt(categoryId, 10, 64)
						if err != nil {
							panic(err)
						}
						posts, err := models.GetPostsByAuthorIdAndCategoryId(atId, ctId)
						if err != nil {
							panic(err)
						}
						context.JSON(http.StatusOK, posts)
					} else if authorId != "" {
						atId, err := strconv.ParseInt(authorId, 10, 64)
						if err != nil {
							panic(err)
						}
						posts, err := models.GetPostsByAuthorId(atId)
						if err != nil {
							panic(err)
						}
						context.JSON(http.StatusOK, posts)
					} else {
						ctId, err := strconv.ParseInt(categoryId, 10, 64)
						if err != nil {
							panic(err)
						}
						posts, err := models.GetPostsByCategoryId(ctId)
						if err != nil {
							panic(err)
						}
						context.JSON(http.StatusOK, posts)
					}
				}
			}
		})

		postRouter.PUT("/:id", func(context *gin.Context) {
			id := context.Param("id")
			if id != "" {
				postId, err := strconv.ParseInt(id, 10, 64)
				if err != nil {
					panic(err)
				}
				originalPost, err := models.GetOnePostById(postId)
				if err != nil {
					panic(err)
				}
				cId := context.PostForm("categoryId")
				categoryId, err := strconv.ParseInt(cId, 10, 64)
				category, err := models.GetOneCategoryById(categoryId)
				if err != nil {
					panic(err)
				}
				title := context.PostForm("title")
				content := context.PostForm("content")
				basicPost := models.BasicPost{
					Title:        title,
					Content:      content,
					AuthorId:     originalPost.AuthorId,
					AuthorName:   originalPost.AuthorName,
					CategoryId:   categoryId,
					CategoryName: category.Title,
				}
				post, err := models.UpdateOnePost(postId, basicPost)
				if err != nil {
					panic(err)
				}
				context.JSON(http.StatusOK, post)
			} else {
				context.AbortWithError(http.StatusBadRequest, errors.New("there is no post ID param"))
			}
		})

		postRouter.DELETE("/:id", func(context *gin.Context) {
			id := context.Param("id")
			if id != "" {
				postId, err := strconv.ParseInt(id, 10, 64)
				if err != nil {
					panic(err)
				}
				post, err := models.DeleteOnePost(postId)
				if err != nil {
					panic(err)
				}
				context.JSON(http.StatusOK, post)
			} else {
				context.AbortWithError(http.StatusBadRequest, errors.New("there is no post ID param"))
			}
		})
	}
}
