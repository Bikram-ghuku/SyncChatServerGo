package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func extractBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("bad header value given")
	}

	jwtToken := strings.Split(header, " ")
	if len(jwtToken) != 2 {
		return "", errors.New("incorrectly formatted authorization header")
	}

	return jwtToken[1], nil
}

func JWTTokenCheck(c *gin.Context) {
	jwtToken, err := extractBearerToken(c.GetHeader("Authorization"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err})
	}

	token, err := jwt.Parse(jwtToken, func(t *jwt.Token) (interface{}, error) {
		if _, OK := t.Method.(*jwt.SigningMethodHMAC); !OK {
			return nil, errors.New("bad signed method received")
		}

		return []byte(os.Getenv("JWT_KEY")), nil
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Bad jwt token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		fmt.Println(claims)
	} else {
		fmt.Printf("Invalid JWT Token")
	}
	c.Set("data", claims)
	c.Next()
}
