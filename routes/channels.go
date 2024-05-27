package routes

import (
	"github.com/Bikram-ghuku/SyncChatServerGo/controllers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupChannelsRoutes(DB *gorm.DB, superRoute *gin.RouterGroup){
  userRoutes := superRoute.Group("/channels")
  {
    userRoutes.GET("/", func(ctx *gin.Context) {
      controllers.GetChannels(ctx, DB);
    })
  }
}
