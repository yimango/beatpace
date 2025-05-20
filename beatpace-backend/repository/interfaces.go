package repository

import (
	"context"

	"github.com/yimango/beatpace-backend/model"
)

// UserRepository handles user data storage operations
type UserRepository interface {
	GetUser(ctx context.Context, id string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
	UpdateUser(ctx context.Context, user *model.User) error
	GetUserBySpotifyID(ctx context.Context, spotifyID string) (*model.User, error)
}

// TokenRepository handles session and token storage operations
type TokenRepository interface {
	SaveSession(ctx context.Context, session *model.Session) error
	GetSessionByToken(ctx context.Context, token string) (*model.Session, error)
	DeleteSession(ctx context.Context, token string) error
	SaveSpotifyToken(ctx context.Context, token *model.SpotifyToken) error
	GetSpotifyToken(ctx context.Context, userID string) (*model.SpotifyToken, error)
} 