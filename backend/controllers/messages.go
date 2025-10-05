package controllers

import (
	"net/http"

	"github.com/Bikram-ghuku/SyncChatServerGo/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MessageResponse struct {
    models.Messages     
    IsSelf bool `json:"self"`
}
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
	result := DB.Where("chat_id = ?",  msgRequest.ChatId).Order("created_at DESC").Find((&findMsgs))

	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No messages found for this chat"})
		return
	}

	response := make([]MessageResponse, len(findMsgs))
	for i, msg := range findMsgs {
		response[i] = MessageResponse{
			Messages: msg,
			IsSelf:   msg.UserId == claims.UserId,
		}
	}


	c.JSON(http.StatusOK, response)

}