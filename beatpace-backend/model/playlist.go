package model

import "time"

// PlaylistResponse is what you return from /api/generate-playlist
type PlaylistResponse struct {
  URL       string    `json:"url"`
  Tracks    []string  `json:"tracks"`
  CreatedAt time.Time `json:"created_at"`
}
