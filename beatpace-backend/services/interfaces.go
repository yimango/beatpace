package services

import (
	"context"

	"github.com/yimango/beatpace-backend/model"
	"github.com/zmb3/spotify/v2"
)

// AuthService handles user authentication and session management
type AuthService interface {
	CreateSession(ctx context.Context, userID string) (*model.Session, error)
	ValidateSession(ctx context.Context, token string) (*model.Session, error)
	RevokeSession(ctx context.Context, token string) error
}

// UserService handles user-related operations
type UserService interface {
	GetUser(ctx context.Context, id string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
	UpdateUser(ctx context.Context, user *model.User) error
	GetUserBySpotifyID(ctx context.Context, spotifyID string) (*model.User, error)
}

// SpotifyService handles Spotify API interactions
type SpotifyService interface {
	GetUserProfile(ctx context.Context, userID string) (*model.SpotifyToken, error)
	RefreshToken(ctx context.Context, userID string) (*model.SpotifyToken, error)
	HandleCallback(ctx context.Context, code string) (string, error)
	GeneratePlaylistForPace(ctx context.Context, userID string, paceSeconds int, gender string, height int) (*PlaylistResponse, error)
	GetClient(ctx context.Context, userID string) *spotify.Client
	GetAuthURL() string
} 