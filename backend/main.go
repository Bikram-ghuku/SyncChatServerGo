package main

import (
	"github.com/Bikram-ghuku/SyncChatServerGo/routes"
	"github.com/Bikram-ghuku/SyncChatServerGo/services"
	"github.com/gin-gonic/gin"
)

func main() {

	DB := services.InitDB()
	app := gin.New()
	router := app.Group("/")
	routes.SetupChannelsRoutes(DB, router)
	routes.SetupUserRoutes(DB, router)

	app.Run()
}
