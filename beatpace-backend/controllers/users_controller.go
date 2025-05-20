// controllers/users_controllers.go
package controllers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yimango/beatpace-backend/model"
	"github.com/yimango/beatpace-backend/services"
	"github.com/yimango/beatpace-backend/types"
)

type UserController struct {
	userService    services.UserService
	authService    services.AuthService
	spotifyService services.SpotifyService
}

func NewUserController(
	userService services.UserService,
	authService services.AuthService,
	spotifyService services.SpotifyService,
) *UserController {
	return &UserController{
		userService:    userService,
		authService:    authService,
		spotifyService: spotifyService,
	}
}

// Register handles user registration
func (uc *UserController) Register(c *gin.Context) {
	var req types.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &model.User{
		ID:           uuid.New(),
		SpotifyUserID: req.SpotifyUserID,
		Email:        req.Email,
	}

	if err := uc.userService.CreateUser(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	session, err := uc.authService.CreateSession(c.Request.Context(), user.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user":    user,
		"token":   session.Token,
		"expires": session.ExpiresAt,
	})
}

// Login handles user login
func (uc *UserController) Login(c *gin.Context) {
	var req types.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := uc.userService.GetUserBySpotifyID(c.Request.Context(), req.SpotifyUserID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	session, err := uc.authService.CreateSession(c.Request.Context(), user.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":    user,
		"token":   session.Token,
		"expires": session.ExpiresAt,
	})
}

// MeHandler returns the current user's information
func (uc *UserController) MeHandler(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := uc.userService.GetUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get Spotify profile information
	spotifyToken, err := uc.spotifyService.GetUserProfile(c.Request.Context(), userID)
	if err != nil {
		// Log the error but don't fail the request
		log.Printf("Failed to get Spotify profile: %v", err)
	}

	response := gin.H{"user": user}
	if spotifyToken != nil {
		response["spotify_token"] = spotifyToken
	}

	c.JSON(http.StatusOK, response)
}

// Callback handles the Spotify OAuth callback
func (uc *UserController) Callback(c *gin.Context) {
	// 1) Read code
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing code"})
		return
	}

	// 2) Handle the callback through our service
	userID, err := uc.spotifyService.HandleCallback(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to handle callback"})
		return
	}

	// 3) Create a session
	session, err := uc.authService.CreateSession(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create session"})
		return
	}

	// 4) Redirect to frontend with token
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}
	params := url.Values{}
	params.Add("token", session.Token)
	params.Add("expires", session.ExpiresAt.String())

	redirectURL := fmt.Sprintf("%s?%s", frontendURL, params.Encode())
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// SignOut handles user sign out
func (uc *UserController) SignOut(c *gin.Context) {
	// Get the token from the Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing authorization header"})
		return
	}

	// Extract token from "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid authorization header"})
		return
	}

	token := parts[1]

	// Revoke the session
	if err := uc.authService.RevokeSession(c.Request.Context(), token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to revoke session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully signed out"})
}
