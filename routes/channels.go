package routes

import (
	"github.com/Bikram-ghuku/SyncChatServerGo/controllers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupChannelsRoutes(DB *gorm.DB)*gin.Engine{
  r := gin.Default()
  r.POST("/channels", func(ctx *gin.Context) {
    controllers.GetChannels(ctx, DB);
  })

  return r;
}
