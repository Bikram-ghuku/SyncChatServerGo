package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type userId struct {
	User string `json:"user"`
}

func GetChannels(c *gin.Context, DB *gorm.DB) {
	senderData := userId{}

	if err := c.BindJSON(&senderData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "username and password is needed"})
		panic("Get channels error")
	}
	c.Bind(&senderData)
	fmt.Println(senderData.User)
	//DB.First(&models.Users{}, senderData)
	c.JSON(http.StatusOK, gin.H{"data": senderData.User})
}
