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

type CheckStruct struct {
	Name       string `json:"name" gorm:"column:name"`
	UserId     string `json:"userId" gorm:"column:user_id"`
	Url        string `json:"url" gorm:"column:url"`
	LastOnline string `json:"lastOnline" gorm:"column:last_online"`
	ChatId     string `json:"chanId" gorm:"column:chat_id"`
	LastMsg    string `json:"lastMsg" gorm:"column:last_msg"`
	NumUnread  string `json:"noUnread" gorm:"column:num_unread"`
}

func getJwtData(c *gin.Context) ClaimStruct {
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

	return claimStruct
}

func GetChannels(c *gin.Context, DB *gorm.DB) {
	//DB.First(&models.Users{}, senderData)

	claimsStruct := getJwtData(c)

	checkData := []CheckStruct{}

	query := `SELECT 
		a.name, 
		a.user_id, 
		a.url as url, 
		a.last_online, 
		b.chat_id, 
		b.last_msg, 
		0 as num_unread
		FROM users a 
		INNER JOIN chats b 
		ON b.sender_id = a.user_id 
		WHERE b.chat_id IN (SELECT chat_id FROM chats WHERE sender_id = ?) 
		AND b.sender_id != ?`

	DB.Raw(query, claimsStruct.UserId, claimsStruct.UserId).Scan(&checkData)

	c.JSON(http.StatusOK, checkData)
}

func AddChannels(c *gin.Context, DB *gorm.DB) {
	type ReqData struct {
		Email string `json:"email"`
	}
	claimsStruct := getJwtData(c)

	reqData := ReqData{}

	if err := c.BindJSON(&reqData); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "error in request"})
		return
	}

	findUser := models.Users{}
	result := DB.Find(&findUser, "email = ?", reqData.Email)

	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}

	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "No user with the given Email found"})
		return
	}

	type CheckStruct struct {
		ChatId uuid.UUID `json:"chat_id"`
		SendId uuid.UUID `json:"sender_id"`
		RecvId uuid.UUID `json:"receiver_id"`
	}

	query := fmt.Sprintf("SELECT DISTINCT a.chat_id as chat_id, a.sender_id as sender_id, b.sender_id as receiver_id FROM chats a INNER JOIN chats b ON a.chat_id = b.chat_id AND a.sender_id != b.sender_id WHERE a.sender_id = '%s' AND b.sender_id = '%s' ", claimsStruct.UserId, findUser.UserId)

	checkData := []CheckStruct{}

	DB.Raw(query).Scan(&checkData)

	checkLen := len(checkData)

	if checkLen != 0 {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": "Channel between the users exsists"})
		return
	}

	newChat_a := models.Chats{
		SenderId: claimsStruct.UserId,
	}

	result = DB.Create(&newChat_a)

	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Error creating new Chat"})
		return
	}

	newChat_b := models.Chats{
		SenderId: findUser.UserId,
		ChatId:   newChat_a.ChatId,
	}

	result = DB.Create(&newChat_b)

	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Error creating new Chat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chat successfully created"})
}
