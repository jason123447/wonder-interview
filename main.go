package main

import (
	"wonder-interview/internal/config"
	"wonder-interview/internal/handlers"
	"wonder-interview/internal/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitConfig()
	utils.InitLogger()
	r := gin.Default()
	r.Use(gin.CustomRecovery(func(c *gin.Context, err any) {
		utils.Logger.Error(err)
	}))

	r.POST("/login", handlers.LoginHandler)

	ws := r.Group("/ws")
	ws.Use(utils.AuthenticateJWTMiddleware(config.SECRET_KEY))
	ws.GET("/:channelID", handlers.SocketHandler)
	r.Run(":8080")
}
