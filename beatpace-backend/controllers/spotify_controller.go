package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yimango/beatpace-backend/services"
)

// GeneratePlaylist handles the playlist generation request
func GeneratePlaylist(c *gin.Context) {
	var request struct {
		AuthCode      string  `json:"authCode"`
		PaceUnit      string  `json:"paceUnit"`
		PaceInSeconds int     `json:"paceInSeconds"`
		Gender        string  `json:"gender"`
		Height        float64 `json:"height"`
		HeightUnit    string  `json:"heightUnit"`
	}

	// Parse JSON from the request
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Exchange the authorization code for an access token
	response, err := services.GetSpotifyAccessToken(request.AuthCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get Spotify access token"})
		return
	}

	// Call the service to generate the playlist
	playlistLink, err := services.CreatePlaylist((*response).AccessToken, request.PaceUnit, request.PaceInSeconds, request.Gender, request.Height, request.HeightUnit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate playlist"})
		return
	}

	// Return the playlist link to the frontend
	c.JSON(http.StatusOK, gin.H{"playlist_link": playlistLink})
}
