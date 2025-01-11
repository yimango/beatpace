package controllers

import (
	"net/http"
	"github.com/yimango/beatpace-backend/services"

	"github.com/gin-gonic/gin"
)

// GeneratePlaylist handles the playlist generation request
func GeneratePlaylist(c *gin.Context) {
	var request struct {
		AccessToken  string  `json:"accessToken"`
		PaceUnit     string  `json:"paceUnit"`
		PaceInSeconds int     `json:"paceInSeconds"`
		Gender       string  `json:"gender"`
		Height       float64 `json:"height"`
		HeightUnit   string  `json:"heightUnit"`
	}

	// Parse JSON from the request
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Call the service to generate the playlist
	playlistLink, err := services.CreatePlaylist(request.AccessToken, request.PaceUnit, request.PaceInSeconds, request.Gender, request.Height, request.HeightUnit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate playlist"})
		return
	}

	// Return the playlist link to the frontend
	c.JSON(http.StatusOK, gin.H{"playlist_link": playlistLink})
}
