package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Bikram-ghuku/SyncChatServerGo/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClaimStruct struct {
	Email  string    `json:"email"`
	Exp    int64     `json:"exp"`
	UserId uuid.UUID `json:"userId"`
}

func GetChannels(c *gin.Context, DB *gorm.DB) {
	//DB.First(&models.Users{}, senderData)
	data, present := c.Get("data")
	if !present {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "JWT authorisation failed"})
	}
	data_json, err := json.Marshal(data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "JWT data error"})
	}

	claimStruct := ClaimStruct{}

	err = json.Unmarshal(data_json, &claimStruct)

	if err != nil {
		panic(err)
	}

	findChannels := models.Chats{}
	result := DB.Find(&findChannels, "sender_id = ?", claimStruct.UserId)

	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	fmt.Println(result.RowsAffected)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{"data": gin.H{}})
	} else {
		c.JSON(http.StatusOK, gin.H{"data": findChannels})

	}
}
