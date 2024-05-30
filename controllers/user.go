package controllers

import (
	"net/http"

	"github.com/Bikram-ghuku/SyncChatServerGo/models"
	"github.com/Bikram-ghuku/SyncChatServerGo/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type signUp struct {
	Email string `form:"email"`
	Name  string `form:"name"`
	Pswd  string `form:"pswd"`
}

func Register(c *gin.Context, DB *gorm.DB) {
	var regUser signUp
	err := c.BindJSON(&regUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error in request"})
		return
	}
	var findUser = models.Users{}
	result := DB.Find(&findUser, "email = ?", regUser.Email)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
	}
	if result.RowsAffected != 0 {
		c.JSON(http.StatusConflict, gin.H{"message": "User already exists"})
		return
	}

	hashPswd, err := services.HashPassword(regUser.Pswd)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error hashing password"})
		return
	}

	var newUser = models.Users{
		Email:    regUser.Email,
		Password: hashPswd,
		Name:     regUser.Name,
	}

	resData := DB.Create(newUser)

	if resData.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating new user"})
		return
	}
}
