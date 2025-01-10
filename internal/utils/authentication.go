package utils

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"
	"wonder-interview/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func GenerateJWT(userID int, secretKey string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func AuthenticateJWTMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		channelIDStr := c.Param("channelID")
		userID, err := strconv.Atoi(channelIDStr)
		if err != nil {
			c.Abort()
			panic(NewErrorResponse(http.StatusBadRequest, "Invalid user_id", err))
		}

		user := models.FindUserByID(userID)
		if user == nil {
			c.Abort()
			panic(NewErrorResponse(http.StatusUnauthorized, "Invalid user_id", err))
		}
		// log.Println("Current User: ", models.MockUsers[idx].Account)
		tokenStr := c.Query("token")
		if tokenStr == "" {
			c.Abort()
			panic(NewErrorResponse(http.StatusUnauthorized, "Invalid user_id", err))
		}
		// log.Println("Token: ", token)
		token, err := jwt.Parse(
			tokenStr,
			func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}
				return []byte(secretKey), nil
			},
		)
		if err != nil || !token.Valid {
			c.Abort()
			panic(NewErrorResponse(http.StatusUnauthorized, "Invalid token", err))
		}
		userIDFromToken := int(token.Claims.(jwt.MapClaims)["user_id"].(float64))
		log.Println("Channel ID: ", channelIDStr)
		log.Println("User ID from token: ", userIDFromToken)
		if userIDFromToken != userID {
			c.Abort()
			panic(NewErrorResponse(http.StatusUnauthorized, "Invalid websocket channel", err))
		}
	}
}
