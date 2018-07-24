package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/wtzeng1/golb/models"
	"net/http"
	"strconv"
)

func RegisterCategoryController() {
	categoryController := router.Group("/category")

	{
		categoryController.POST("/", func(context *gin.Context) {
			title := context.PostForm("title")
			authorId := context.PostForm("authorId")
			aId, err := strconv.ParseInt(authorId, 10, 64)
			if err != nil {
				panic(err)
			}
			author, err := models.GetAuthorById(aId)
			if err != nil {
				panic(err)
			}
			authorName := author.Name
			basicCategory := models.BasicCategory{
				Title:      title,
				AuthorId:   aId,
				AuthorName: authorName,
			}
			category, err := models.AddOneCategory(basicCategory)
			if err != nil {
				panic(err)
			}
			context.JSON(http.StatusOK, category)
		})

		categoryController.GET("/", func(context *gin.Context) {
			id := context.Query("id")
			authorId := context.Query("authorId")
			title := context.Query("title")
			if id == "" && authorId == "" && title == "" {
				categories, err := models.GetAllCategories()
				if err != nil {
					panic(err)
				}
				context.JSON(http.StatusOK, categories)
			} else if id != "" {
				categoryId, err := strconv.ParseInt(id, 10, 64)
				if err != nil {
					panic(err)
				}
				category, err := models.GetOneCategoryById(categoryId)
				if err != nil {
					panic(err)
				}
				context.JSON(http.StatusOK, category)
			} else {
				if authorId != "" && title != "" {
					aId, err := strconv.ParseInt(authorId, 10, 64)
					if err != nil {
						panic(err)
					}
					categories, err := models.GetCategoriesByAuthorIdAndTitle(aId, title)
					if err != nil {
						panic(err)
					}
					context.JSON(http.StatusOK, categories)
				} else if authorId != "" {
					aId, err := strconv.ParseInt(authorId, 10, 64)
					if err != nil {
						panic(err)
					}
					categories, err := models.GetCategoriesByAuthorId(aId)
					if err != nil {
						panic(err)
					}
					context.JSON(http.StatusOK, categories)
				} else {
					categories, err := models.GetCategoriesByTitle(title)
					if err != nil {
						panic(err)
					}
					context.JSON(http.StatusOK, categories)
				}
			}
		})

		categoryController.PUT("/:id", func(context *gin.Context) {
			id := context.Param("id")
			if id != "" {
				cId, err := strconv.ParseInt(id, 10, 64)
				if err != nil {
					panic(err)
				}
				title := context.PostForm("title")
				category, err := models.UpdateOneCategory(cId, title)
				if err != nil {
					panic(err)
				}
				context.JSON(http.StatusOK, category)
			} else {
				context.AbortWithError(http.StatusBadRequest, errors.New("there is no category ID param"))
			}
		})

		categoryController.DELETE("/:id", func(context *gin.Context) {
			id := context.Param("id")
			if id != "" {
				cId, err := strconv.ParseInt(id, 10, 64)
				if err != nil {
					panic(err)
				}
				category, err := models.DeleteOneCategory(cId)
				if err != nil {
					panic(err)
				}
				context.JSON(http.StatusOK, category)
			} else {
				context.AbortWithError(http.StatusBadRequest, errors.New("there is no category ID param"))
			}
		})
	}
}
