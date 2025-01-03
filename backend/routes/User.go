package routes

import (
	"github.com/Bikram-ghuku/SyncChatServerGo/controllers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupUserRoutes(DB *gorm.DB, superRoute *gin.RouterGroup) {
	userRoutes := superRoute.Group("/users")
	{
		userRoutes.POST("/register", func(ctx *gin.Context) {
			controllers.Register(ctx, DB)
		})

		userRoutes.POST("/login", func(ctx *gin.Context) {
			controllers.Login(ctx, DB)
		})

		userRoutes.POST("/ghreg", func(ctx *gin.Context) {
			controllers.GhAuth(ctx, DB)
		})
	}
}
