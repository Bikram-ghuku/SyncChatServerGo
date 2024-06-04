package routes

import (
	"github.com/Bikram-ghuku/SyncChatServerGo/controllers"
	"github.com/Bikram-ghuku/SyncChatServerGo/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupChannelsRoutes(DB *gorm.DB, superRoute *gin.RouterGroup) {
	channelRoutes := superRoute.Group("/channels")
	{
		channelRoutes.Use(middleware.JWTTokenCheck)
		channelRoutes.GET("/channels", func(ctx *gin.Context) {
			controllers.GetChannels(ctx, DB)
		})
	}
}
