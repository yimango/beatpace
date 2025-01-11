package services

import (
	"fmt"
)

// CreatePlaylist generates a Spotify playlist based on the user's input
func CreatePlaylist(accessToken, paceUnit string, paceInSeconds int, gender string, height float64, heightUnit string) (string, error) {
	// Convert height to centimeters if needed
	if heightUnit == "in" {
		height = height * 2.54
	}

	// Calculate the desired BPM based on pace
	pacePerKm := float64(paceInSeconds) / 60.0
	bpm := int(180.0 / pacePerKm)

	// Call Spotify API to create a playlist
	playlistID, err := createSpotifyPlaylist(accessToken, bpm)
	if err != nil {
		return "", err
	}

	// Return the Spotify playlist link
	return fmt.Sprintf("https://open.spotify.com/playlist/%s", playlistID), nil
}

// Helper function to interact with the Spotify API
func createSpotifyPlaylist(accessToken string, bpm int) (string, error) {
	// Example Spotify API call to create a playlist
	// You need to implement the actual API call with error handling
	playlistID := "examplePlaylistID"
	return playlistID, nil
}
