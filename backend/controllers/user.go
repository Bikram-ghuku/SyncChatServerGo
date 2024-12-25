package controllers

import (
	"encoding/json"
	"fmt"
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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "error in request"})
		return
	}

	findUser := models.Users{}
	result := DB.Find(&findUser, "email = ?", regUser.Email)

	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}

	if result.RowsAffected != 0 {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": "User already exists"})
		return
	}

	hashPswd, err := services.HashPassword(regUser.Pswd)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Error hashing password"})
		return
	}

	var newUser = models.Users{
		Email:    regUser.Email,
		Password: hashPswd,
		Name:     regUser.Name,
	}

	resData := DB.Create(&newUser)

	if resData.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Error creating new user"})
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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "error in request"})
		return
	}

	findUser := models.Users{}
	result := DB.Find(&findUser, "email = ?", loginUser.Email)

	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}

	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "No user with the email is not found"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(findUser.Password), []byte(loginUser.Pswd))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "username or password is wrong"})
		return
	}

	signKey := os.Getenv("JWT_KEY")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"email": findUser.Email, "exp": time.Now().Add(time.Hour * 24).Unix(), "userId": findUser.UserId})

	tokenString, err := token.SignedString([]byte(signKey))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "JWT Signing error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login Successful", "token": tokenString, "name": findUser.Name, "profile_url": findUser.Url})

}

func GhAuth(c *gin.Context, DB *gorm.DB) {
	type BodyReg struct {
		GhCode string `json:"code"`
	}

	type GithubAccessTokenResponse struct {
		AccessToken string `json:"access_token"`
		Scope       string `json:"scope"`
		TokenType   string `json:"token_type"`
	}

	type GithubUserResponse struct {
		Name      string `json:"name"`
		Login     string `json:"login"`
		ID        int    `json:"id"`
		AvatarURL string `json:"avatar_url"`
	}

	bodyReg := BodyReg{}
	if err := c.Bind(&bodyReg); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "error in request"})
		return
	}

	gh_pubKey := os.Getenv("GH_CLIENT_ID")
	gh_pvtKey := os.Getenv("GH_PRIVATE_ID")
	uri := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", gh_pubKey, gh_pvtKey, bodyReg.GhCode)

	req, _ := http.NewRequest("POST", uri, nil)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()
	var tokenResponse GithubAccessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	req, _ = http.NewRequest("GET", "https://api.github.com/user", nil)
	req.Header.Set("Authorization", "Bearer "+tokenResponse.AccessToken)

	resp, err = client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	var userResponse GithubUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	name := userResponse.Name
	if name == "" {
		name = userResponse.Login
	}

	uname := userResponse.Login
	id := fmt.Sprintf("%d", userResponse.ID)
	url := userResponse.AvatarURL

	storeUser := models.Users{
		Email:    uname,
		Name:     name,
		Url:      url,
		Password: id,
	}
	if id != "0" {
		DB.Create(&storeUser)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid credentials"})
		return
	}

	findUser := models.Users{
		Email: uname,
	}

	result := DB.First(&findUser)

	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid credentials"})
		return
	}

	if findUser.Password != id || result.RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": findUser.Email,
		"name":  findUser.Name,
		"id":    findUser.Email,
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("TOKEN")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"name":  findUser.Name,
		"email": findUser.Email,
		"id":    findUser.Email,
		"url":   findUser.Url,
	})
}
