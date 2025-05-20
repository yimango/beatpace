// controllers/controller.go
package controllers

import (
	"net/http"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yimango/beatpace-backend/middleware"
	"github.com/yimango/beatpace-backend/model"
	"github.com/yimango/beatpace-backend/services"
)

// Controller handles HTTP requests
type Controller struct {
	userService    services.UserService
	authService    services.AuthService
	spotifyService services.SpotifyService
}

// NewController creates a new controller instance
func NewController(
	userService services.UserService,
	authService services.AuthService,
	spotifyService services.SpotifyService,
) *Controller {
	return &Controller{
		userService:    userService,
		authService:    authService,
		spotifyService: spotifyService,
	}
}

type RegisterRequest struct {
	SpotifyUserID string  `json:"spotify_user_id" binding:"required"`
	Email         *string `json:"email"`
}

type LoginRequest struct {
	SpotifyUserID string `json:"spotify_user_id" binding:"required"`
}

type GeneratePlaylistRequest struct {
	PaceSeconds int    `json:"pace_seconds" binding:"required"`
	Gender      string `json:"gender" binding:"required"`
	Height      int    `json:"height" binding:"required"`
}

// Register handles user registration
func (c *Controller) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &model.User{
		ID:           uuid.New(),
		SpotifyUserID: req.SpotifyUserID,
		Email:        req.Email,
	}

	if err := c.userService.CreateUser(ctx, user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	session, err := c.authService.CreateSession(ctx, user.ID.String())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"user":    user,
		"token":   session.Token,
		"expires": session.ExpiresAt,
	})
}

// Login handles user login
func (c *Controller) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.userService.GetUserBySpotifyID(ctx, req.SpotifyUserID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	session, err := c.authService.CreateSession(ctx, user.ID.String())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user":    user,
		"token":   session.Token,
		"expires": session.ExpiresAt,
	})
}

// MeHandler returns the current user's information
func (c *Controller) MeHandler(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := c.userService.GetUser(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get Spotify profile information
	spotifyToken, err := c.spotifyService.GetUserProfile(ctx, userID)
	if err != nil {
		// Log the error but don't fail the request
		log.Printf("Failed to get Spotify profile: %v", err)
	}

	response := gin.H{"user": user}
	if spotifyToken != nil {
		response["spotify_token"] = spotifyToken
	}

	ctx.JSON(http.StatusOK, response)
}

// Callback handles the Spotify OAuth callback
func (c *Controller) Callback(ctx *gin.Context) {
	code := ctx.Query("code")
	if code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing code parameter"})
		return
	}

	userID, err := c.spotifyService.HandleCallback(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	session, err := c.authService.CreateSession(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token":   session.Token,
		"expires": session.ExpiresAt,
	})
}

// GeneratePlaylist creates a playlist based on user parameters
func (c *Controller) GeneratePlaylist(ctx *gin.Context) {
	var req GeneratePlaylistRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := ctx.GetString("userID")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	playlist, err := c.spotifyService.GeneratePlaylistForPace(ctx, userID, req.PaceSeconds, req.Gender, req.Height)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, playlist)
}

// RegisterRoutes registers all the controller routes
func (c *Controller) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	
	// Public routes
	api.POST("/login", c.Login)
	api.POST("/register", c.Register)
	api.GET("/callback", c.Callback)

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.JWT(c.authService))
	{
		protected.GET("/me", c.MeHandler)
		protected.POST("/generate-playlist", c.GeneratePlaylist)
	}
}