package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/yimango/beatpace-backend/controllers"
	"github.com/yimango/beatpace-backend/db"
	"github.com/yimango/beatpace-backend/middleware"
	"github.com/yimango/beatpace-backend/repository"
	"github.com/yimango/beatpace-backend/services"
)

func main() {
	// Set Gin to release mode
	gin.SetMode(gin.ReleaseMode)

	// Debug: Print all environment variables (excluding sensitive ones)
	log.Printf("Environment variables loaded:")
	log.Printf("JWT_SECRET exists: %v", os.Getenv("JWT_SECRET") != "")
	log.Printf("CLIENT_ID exists: %v", os.Getenv("CLIENT_ID") != "")
	log.Printf("CLIENT_SECRET exists: %v", os.Getenv("CLIENT_SECRET") != "")
	log.Printf("REDIRECT_URI exists: %v", os.Getenv("REDIRECT_URI") != "")

	// Verify required environment variables
	requiredEnvVars := []string{"JWT_SECRET", "CLIENT_ID", "CLIENT_SECRET", "REDIRECT_URI"}
	for _, envVar := range requiredEnvVars {
		value := os.Getenv(envVar)
		if value == "" {
			log.Fatalf("Required environment variable %s is not set", envVar)
		}
		log.Printf("Found %s: %s", envVar, value[:min(5, len(value))]+"...")
	}

	// 1) initialize the DB connection and run migrations
	sqlDB, err := db.InitDatabase()
	if err != nil {
		log.Fatalf("failed to initialize DB: %v", err)
	}
	defer sqlDB.Close()

	// 2) wire up your repositories
	userRepo := repository.NewUserRepo(sqlDB)
	tokenRepo := repository.NewTokenRepo(sqlDB)

	// Print Spotify configuration for debugging
	log.Printf("Spotify Configuration - Client ID: %s, Redirect URI: %s", os.Getenv("CLIENT_ID"), os.Getenv("REDIRECT_URI"))

	// 3) wire up your services
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo, tokenRepo)
	spotifyService := services.NewSpotifyService(userRepo, tokenRepo)

	// 4) create your controllers
	userController := controllers.NewUserController(userService, authService, spotifyService)
	spotifyController := controllers.NewSpotifyController(spotifyService)

	// 5) create the Gin router
	router := gin.Default()

	// Configure trusted proxies
	router.SetTrustedProxies([]string{"127.0.0.1"})

	// 6) configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 7) register routes
	api := router.Group("/api")
	{
		// Public routes
		api.POST("/login", userController.Login)
		api.POST("/register", userController.Register)
		api.GET("/callback", userController.Callback)
		api.GET("/spotify-auth", spotifyController.GetAuthURL)

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.JWT(authService))
		{
			protected.GET("/me", userController.MeHandler)
			protected.POST("/generate-playlist", spotifyController.GeneratePlaylist)
			protected.POST("/signout", userController.SignOut)
		}
	}

	// 8) start the server
	log.Println("Server running on http://localhost:3001")
	router.Run(":3001")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
