package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/wtzeng1/golb/models"
	"net/http"
	"strconv"
)

func RegisterAuthorController() {
	authorRouter := router.Group("/author")

	{
		authorRouter.POST("/", func(context *gin.Context) {
			username := context.PostForm("username")
			email := context.PostForm("email")
			password := context.PostForm("password")
			name := context.PostForm("name")
			author := models.BasicAuthor{
				Username: username,
				Email:    email,
				Password: password,
				Name:     name,
			}
			a, err := models.AddOneAuthor(author)
			if err != nil {
				panic(err)
			}
			context.JSON(http.StatusOK, a)
		})

		authorRouter.GET("/", func(context *gin.Context) {
			id := context.Query("id")
			email := context.Query("email")
			if id != "" {
				id, err := strconv.ParseInt(id, 10, 64)
				if err != nil {
					panic(err)
				}
				author, err := models.GetAuthorById(id)
				if err != nil {
					panic(err)
				}
				context.JSON(http.StatusOK, author)
			} else if email != "" {
				author, err := models.GetAuthorByEmail(email)
				if err != nil {
					panic(err)
				}
				context.JSON(http.StatusOK, author)
			} else {
				authors, err := models.GetAllAuthors()
				if err != nil {
					panic(err)
				}
				context.JSON(http.StatusOK, authors)
			}
		})

		authorRouter.PUT("/:id", func(context *gin.Context) {
			id := context.Param("id")
			if id != "" {
				id, err := strconv.ParseInt(id, 10, 64)
				if err != nil {
					panic(err)
				}
				username := context.PostForm("username")
				email := context.PostForm("email")
				password := context.PostForm("password")
				name := context.PostForm("name")
				author := models.BasicAuthor{
					Username: username,
					Email:    email,
					Password: password,
					Name:     name,
				}
				a, err := models.UpdateOneAuthor(id, author)
				if err != nil {
					panic(err)
				}
				context.JSON(http.StatusOK, a)
			} else {
				context.AbortWithError(http.StatusBadRequest, errors.New("there is no author ID param"))
			}
		})

		authorRouter.DELETE("/:id", func(context *gin.Context) {
			id := context.Param("id")
			if id != "" {
				aId, err := strconv.ParseInt(id, 10, 64)
				if err != nil {
					panic(err)
				}
				author, err := models.DeleteOneAuthor(aId)
				if err != nil {
					panic(err)
				}
				context.JSON(http.StatusOK, author)
			} else {
				context.AbortWithError(http.StatusBadRequest, errors.New("there is no author ID param"))
			}
		})
	}
}
