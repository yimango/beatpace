// main.go
package main

import (
    "log"
    "time"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "github.com/yimango/beatpace-backend/controllers"
    "github.com/yimango/beatpace-backend/middleware"
)

func main() {
    router := gin.Default()

    // 1) CORS must allow credentials from your frontend
    router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"},
        AllowMethods:     []string{"GET", "POST"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

    // 2) OAuth callback (no auth)
    router.GET("/api/callback", controllers.Callback)

    // 3) Protected routes
    api := router.Group("/api")
    api.Use(middleware.JWT())
    {
        api.GET("/me", controllers.MeHandler)
        // ... other routes ...
    }

    log.Println("Server running on http://localhost:3001")
    router.Run(":3001")
}
