package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/yimango/beatpace-backend/model"
)

type tokenRepository struct {
	db *sql.DB
}

func NewTokenRepo(db *sql.DB) *tokenRepository {
	return &tokenRepository{db: db}
}

func (r *tokenRepository) SaveSession(ctx context.Context, session *model.Session) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO sessions (id, user_id, token, created_at, expires_at) VALUES (?, ?, ?, ?, ?)",
		session.ID, session.UserID, session.Token, session.CreatedAt, session.ExpiresAt)
	if err != nil {
		return fmt.Errorf("error saving session: %v", err)
	}
	return nil
}

func (r *tokenRepository) GetSessionByToken(ctx context.Context, token string) (*model.Session, error) {
	fmt.Printf("Getting session by token: %s\n", token)
	var session model.Session
	err := r.db.QueryRowContext(ctx,
		"SELECT id, user_id, token, created_at, expires_at FROM sessions WHERE token = ?",
		token).Scan(&session.ID, &session.UserID, &session.Token, &session.CreatedAt, &session.ExpiresAt)
	if err == sql.ErrNoRows {
		fmt.Println("No session found for token")
		return nil, fmt.Errorf("session not found")
	}
	if err != nil {
		fmt.Printf("Error getting session: %v\n", err)
		return nil, fmt.Errorf("error getting session: %v", err)
	}
	fmt.Printf("Found session for user %s\n", session.UserID.String())
	return &session, nil
}

func (r *tokenRepository) DeleteSession(ctx context.Context, token string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM sessions WHERE token = ?", token)
	if err != nil {
		return fmt.Errorf("error deleting session: %v", err)
	}
	return nil
}

func (r *tokenRepository) SaveSpotifyToken(ctx context.Context, token *model.SpotifyToken) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO spotify_tokens (internal_user_id, spotify_user_id, access_token, refresh_token, generated_at, expires_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			access_token = VALUES(access_token),
			refresh_token = VALUES(refresh_token),
			generated_at = VALUES(generated_at),
			expires_at = VALUES(expires_at)`,
		token.InternalUserID, token.SpotifyUserID, token.AccessToken, token.RefreshToken, token.GeneratedAt, token.ExpiresAt)
	if err != nil {
		return fmt.Errorf("error saving spotify token: %v", err)
	}
	return nil
}

func (r *tokenRepository) GetSpotifyToken(ctx context.Context, userID string) (*model.SpotifyToken, error) {
	var token model.SpotifyToken
	err := r.db.QueryRowContext(ctx,
		"SELECT internal_user_id, spotify_user_id, access_token, refresh_token, generated_at, expires_at FROM spotify_tokens WHERE internal_user_id = ?",
		userID).Scan(&token.InternalUserID, &token.SpotifyUserID, &token.AccessToken, &token.RefreshToken, &token.GeneratedAt, &token.ExpiresAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("spotify token not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error getting spotify token: %v", err)
	}
	return &token, nil
} 