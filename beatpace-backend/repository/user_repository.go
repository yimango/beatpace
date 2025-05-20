package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/yimango/beatpace-backend/model"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUser(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	err := r.db.QueryRowContext(ctx,
		"SELECT id, spotify_user_id, email, created_at, updated_at, last_login FROM users WHERE id = ?",
		id).Scan(&user.ID, &user.SpotifyUserID, &user.Email, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error getting user: %v", err)
	}
	return &user, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user *model.User) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO users (id, spotify_user_id, email, created_at, updated_at) VALUES (?, ?, ?, NOW(), NOW())",
		user.ID, user.SpotifyUserID, user.Email)
	if err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}
	return nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user *model.User) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE users SET spotify_user_id = ?, email = ?, updated_at = NOW() WHERE id = ?",
		user.SpotifyUserID, user.Email, user.ID)
	if err != nil {
		return fmt.Errorf("error updating user: %v", err)
	}
	return nil
}

func (r *userRepository) GetUserBySpotifyID(ctx context.Context, spotifyID string) (*model.User, error) {
	var user model.User
	err := r.db.QueryRowContext(ctx,
		"SELECT id, spotify_user_id, email, created_at, updated_at, last_login FROM users WHERE spotify_user_id = ?",
		spotifyID).Scan(&user.ID, &user.SpotifyUserID, &user.Email, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error getting user: %v", err)
	}
	return &user, nil
}
