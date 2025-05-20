package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/yimango/beatpace-backend/controllers"
	"github.com/yimango/beatpace-backend/db"
	"github.com/yimango/beatpace-backend/middleware"
	"github.com/yimango/beatpace-backend/repository"
	"github.com/yimango/beatpace-backend/services"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
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
