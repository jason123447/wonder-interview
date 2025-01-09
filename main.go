package main

import (
	"wonder-interview/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	ws := r.Group("/ws")
	ws.GET("/:channelID", handlers.SocketHandler)
	r.Run(":8080")
}
