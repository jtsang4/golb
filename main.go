package main

import (
	"github.com/wtzeng1/golb/controllers"
	"github.com/wtzeng1/golb/models"
)

func main() {
	engine := controllers.CreateApp()
	engine.Run(":9090")
	defer models.CloseDB()
}
