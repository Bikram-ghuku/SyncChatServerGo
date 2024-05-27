package routes

import (
	"github.com/Bikram-ghuku/SyncChatServerGo/controllers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupUserRoutes(DB *gorm.DB, superRoute *gin.RouterGroup){
  userRoutes := superRoute.Group("/user")
  {
    userRoutes.POST("/register", func(ctx *gin.Context) {
      controllers.Register(ctx, DB);
    })
  }
}
