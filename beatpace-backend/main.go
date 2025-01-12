package main

import (
	"fmt"
	"log"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/yimango/beatpace-backend/controllers"
)

func main() {
	// Create a Gin router
	router := gin.Default()

	// Define CORS middleware for Gin
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Change this to your frontend's origin
		AllowMethods:     []string{"POST"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))

	// Set up your routes
	router.POST("/api/generate-playlist", controllers.GeneratePlaylist)

	// Start the server
	fmt.Println("Server is running on http://localhost:3001")
	log.Fatal(router.Run(":3001"))
}
