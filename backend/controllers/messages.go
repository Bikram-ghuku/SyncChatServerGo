package controllers

import (
	"net/http"

	"github.com/Bikram-ghuku/SyncChatServerGo/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetMessages(c *gin.Context, DB *gorm.DB) {
	claims := getJwtData(c)

	var msgRequest struct {
		ChatId string `json:"chatId"`
		Multi  int16  `json:"multi"`
	}


	if err := c.BindJSON(&msgRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "error in request"})
		return
	}

	findMsgs := []models.Messages{}
	result := DB.Find(&findMsgs, "user_id = ? AND chat_id = ?", claims.UserId, msgRequest.ChatId)

	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No messages found for this chat"})
		return
	}


	c.JSON(http.StatusOK, findMsgs)

}