package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/Bikram-ghuku/SyncChatServerGo/models"
	"github.com/Bikram-ghuku/SyncChatServerGo/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(c *gin.Context, DB *gorm.DB) {
	type signUp struct {
		Email string `json:"email"`
		Name  string `json:"name"`
		Pswd  string `json:"pswd"`
	}

	regUser := signUp{}

	if err := c.BindJSON(&regUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error in request"})
		return
	}

	findUser := models.Users{}
	result := DB.Find(&findUser, "email = ?", regUser.Email)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
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

	resData := DB.Create(&newUser)

	if resData.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating new user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})

}

func Login(c *gin.Context, DB *gorm.DB) {
	type Login struct {
		Email string `json:"email"`
		Pswd  string `json:"pswd"`
	}

	loginUser := Login{}

	if err := c.Bind(&loginUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error in request"})
		return
	}

	findUser := models.Users{}
	result := DB.Find(&findUser, "email = ?", loginUser.Email)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "No user with the email is not found"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(findUser.Password), []byte(loginUser.Pswd))

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "username or password is wrong"})
		return
	}

	signKey := os.Getenv("JWT_KEY")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"email": findUser.Email, "exp": time.Now().Add(time.Hour * 24).Unix()})

	tokenString, err := token.SignedString([]byte(signKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "JWT Signing error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login Successful", "jwt": tokenString})

}
