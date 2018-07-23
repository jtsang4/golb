package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/wtzeng1/golb/models"
	"net/http"
)

func RegisterAuthorController() {
	authorRouter := router.Group("/author")

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
}
