package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

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

	fmt.Println(claimStruct.Email)
	fmt.Println(string(data_json))
	c.JSON(http.StatusOK, gin.H{"data": data})
}
