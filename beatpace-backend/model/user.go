package model

import (
  "time"

  "github.com/google/uuid"
)

// User represents a Spotify-authenticated user in the system
type User struct {
  ID            uuid.UUID `db:"id"`              // Primary key
  SpotifyUserID string    `db:"spotify_user_id"` // Spotify user identifier
  Email         *string   `db:"email"`           // Optional email address
  CreatedAt     time.Time `db:"created_at"`      // Account creation timestamp
  UpdatedAt     time.Time `db:"updated_at"`      // Last update timestamp
  LastLogin     *time.Time `db:"last_login"`     // Last login timestamp
}
