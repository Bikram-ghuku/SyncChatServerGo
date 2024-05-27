package main

import (
	"github.com/Bikram-ghuku/SyncChatServerGo/routes"
	"github.com/Bikram-ghuku/SyncChatServerGo/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


func main() {

  var DB *gorm.DB
  DB = services.InitDB();
  app := gin.New()
  router := app.Group("/")
  routes.SetupChannelsRoutes(DB, router)
  routes.SetupUserRoutes(DB, router)

  app.Run()
}
