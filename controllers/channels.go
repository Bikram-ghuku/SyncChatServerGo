package controllers

import (
	"net/http"

	"github.com/Bikram-ghuku/SyncChatServerGo/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type userId struct {
  User string `form:"user"`
}

func GetChannels(c* gin.Context, DB *gorm.DB){
  var senderData userId;
  err := c.BindJSON(&senderData);
  if err != nil {
    panic("Get channels error")
  }

  DB.First(&models.Users{}, senderData);
  c.JSON(http.StatusOK, gin.H{"data": senderData.User})
}
