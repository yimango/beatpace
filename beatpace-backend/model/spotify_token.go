package model

import (
	"time"
)

// SpotifyToken represents OAuth tokens for Spotify API access
type SpotifyToken struct {
	InternalUserID string    `db:"internal_user_id"` // Reference to our user
	SpotifyUserID  string    `db:"spotify_user_id"`  // Spotify's user identifier
	AccessToken    string    `db:"access_token"`     // OAuth access token
	RefreshToken   string    `db:"refresh_token"`    // OAuth refresh token
	GeneratedAt    time.Time `db:"generated_at"`     // When the token was generated
	ExpiresAt      time.Time `db:"expires_at"`       // When the token expires
}
