package controllers

import (
	"fmt"
	"net/http"

	"github.com/Bikram-ghuku/SyncChatServerGo/models"
	"github.com/Bikram-ghuku/SyncChatServerGo/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type signUp struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Pswd  string `json:"pswd"`
}

func Register(c *gin.Context, DB *gorm.DB) {
	regUser := signUp{}

	if err := c.BindJSON(&regUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error in request"})
		panic("error in request")
	}

	hashPswd, err := services.HashPassword(regUser.Pswd)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error hashing password"})
		panic("Error hashing password")
	}

	var newUser = models.Users{
		Email:    regUser.Email,
		Password: hashPswd,
		Name:     regUser.Name,
	}

	fmt.Println("user: ", newUser.Email, "Password: ", newUser.Password, "Name: ", newUser.Name)
	resData := DB.Create(&newUser)

	if resData.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating new user"})
		panic("Error creating new user")
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})

}
