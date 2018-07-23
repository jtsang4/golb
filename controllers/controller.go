package controllers

import "github.com/gin-gonic/gin"

var router *gin.RouterGroup

func CreateApp() *gin.Engine {
	engine := gin.Default()
	router = engine.Group("/api")
	RegisterRouters()
	return engine
}

func RegisterRouters() {
	RegisterAuthorController()
}
