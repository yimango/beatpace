package model

import (
	"time"

	"github.com/google/uuid"
)

// Session represents an active user session
type Session struct {
	ID        uuid.UUID `db:"id"`         // Session unique identifier
	UserID    uuid.UUID `db:"user_id"`    // Reference to the user
	Token     string    `db:"token"`      // Session token
	CreatedAt time.Time `db:"created_at"` // Session creation time
	ExpiresAt time.Time `db:"expires_at"` // Session expiration time
} 