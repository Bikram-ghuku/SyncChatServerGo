package routes

import (
	"github.com/Bikram-ghuku/SyncChatServerGo/controllers"
	"github.com/Bikram-ghuku/SyncChatServerGo/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupMessagesRoutes(DB *gorm.DB, superRoute *gin.RouterGroup) {
	channelRoutes := superRoute.Group("/message")
	{
		channelRoutes.Use(middleware.JWTTokenCheck)
		channelRoutes.POST("/getMsgs", func(ctx *gin.Context) {
			controllers.GetMessages(ctx, DB)
		})
	}
}
