package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yimango/beatpace-backend/services"
	"github.com/yimango/beatpace-backend/types"
)

type SpotifyController struct {
	spotifyService services.SpotifyService
}

func NewSpotifyController(spotifyService services.SpotifyService) *SpotifyController {
	return &SpotifyController{
		spotifyService: spotifyService,
	}
}

// GeneratePlaylist handles the playlist generation request
func (sc *SpotifyController) GeneratePlaylist(c *gin.Context) {
	fmt.Println("Received generate playlist request")
	fmt.Printf("Headers: %+v\n", c.Request.Header)

	var req types.GeneratePlaylistRequest
	if err := c.BindJSON(&req); err != nil {
		fmt.Printf("Failed to bind JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request format: %v", err)})
		return
	}

	fmt.Println("Received request with data:", req)

	// Get userID from context (set by JWT middleware)
	userID := c.GetString("userID")
	if userID == "" {
		fmt.Println("No userID found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	fmt.Printf("Generating playlist for user %s\n", userID)

	// Generate the playlist using our service
	playlist, err := sc.spotifyService.GeneratePlaylistForPace(
		c.Request.Context(),
		userID,
		req.PaceInSeconds,
		req.Gender,
		int(req.Height),
	)
	if err != nil {
		fmt.Printf("Failed to generate playlist: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("Successfully generated playlist: %+v\n", playlist)
	c.JSON(http.StatusOK, playlist)
}

// GetAuthURL returns the Spotify authorization URL
func (c *SpotifyController) GetAuthURL(ctx *gin.Context) {
	authURL := c.spotifyService.GetAuthURL()
	ctx.JSON(http.StatusOK, gin.H{
		"authUrl": authURL,
	})
}
