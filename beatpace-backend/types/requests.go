package types

// RegisterRequest represents the registration request payload
type RegisterRequest struct {
	SpotifyUserID string  `json:"spotify_user_id" binding:"required"`
	Email         *string `json:"email"`
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	SpotifyUserID string `json:"spotify_user_id" binding:"required"`
}

// GeneratePlaylistRequest represents the playlist generation request payload
type GeneratePlaylistRequest struct {
	PaceInSeconds int     `json:"paceInSeconds" binding:"required"`
	Gender        string  `json:"gender" binding:"required"`
	Height        float64 `json:"height" binding:"required"`
} 