package main

import (
	"os"
	"upvotesystem/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8001"
	}
	router := gin.New()
	router.Use(gin.Logger())
	group := router.Group("") // /api/v1
	routes.PostRoutes(group)
	router.Run(":" + port)
}
