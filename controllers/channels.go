package controllers;

import(
  "github.com/gin-gonic/gin"
	"gorm.io/gorm"
  "github.com/Bikram-ghuku/SyncChatServerGo/models"
)

func getChannels(c* gin.Context, DB *gorm.DB){
  var senderData string;
  err := c.BindJSON(&senderData);
  if err != nil {
    panic("Get channels error")
  }

  DB.First(&models.Users{}, senderData);
}
