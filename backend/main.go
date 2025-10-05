package main

import (
	"github.com/Bikram-ghuku/SyncChatServerGo/routes"
	"github.com/Bikram-ghuku/SyncChatServerGo/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	DB := services.InitDB()
	app := gin.New()

	corsRule := cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	app.Use(corsRule)

	router := app.Group("/")
	routes.SetupChannelsRoutes(DB, router)
	routes.SetupUserRoutes(DB, router)
	routes.SetupMessagesRoutes(DB, router)
	app.Run()
}
