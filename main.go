package main

import (
	"wonder-interview/docs"
	"wonder-interview/internal/config"
	"wonder-interview/internal/handlers"
	"wonder-interview/internal/utils"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @Title wonser-interview
// @Version 1.0
// @Description API for wonser-interview
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	config.InitConfig()
	utils.InitLogger()
	r := gin.Default()
	r.Use(gin.CustomRecovery(func(c *gin.Context, err any) {
		utils.Logger.Error(err)
	}))

	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.POST("/login", handlers.LoginHandler)

	ws := r.Group("/ws")
	ws.Use(utils.AuthenticateJWTMiddleware(config.SECRET_KEY))
	ws.GET("/:channelID", handlers.SocketHandler)
	r.Run(":8080")
}
