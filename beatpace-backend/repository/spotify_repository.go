package repository

import (
  "context"
  "database/sql"

  "github.com/yimango/beatpace-backend/model"
)

type SpotifyTokenRepo struct {
  db *sql.DB
}

func NewSpotifyTokenRepo(db *sql.DB) *SpotifyTokenRepo {
  return &SpotifyTokenRepo{db: db}
}

// Upsert inserts or updates the Spotify tokens for a given internal user ID.
func (r *SpotifyTokenRepo) Upsert(ctx context.Context, t *model.SpotifyToken) error {
  const q = `
    INSERT INTO spotify_tokens
      (internal_user_id, spotify_user_id, access_token, refresh_token, generated_at, expires_at)
    VALUES (?, ?, ?, ?, ?, ?)
    ON DUPLICATE KEY UPDATE
      access_token  = VALUES(access_token),
      refresh_token = VALUES(refresh_token),
      generated_at  = VALUES(generated_at),
      expires_at    = VALUES(expires_at)
  `
  _, err := r.db.ExecContext(ctx, q,
    t.InternalUserID,
    t.SpotifyUserID,
    t.AccessToken,
    t.RefreshToken,
    t.GeneratedAt,
    t.ExpiresAt,
  )
  return err
}

// FindByInternalUserID loads a user’s Spotify tokens.
func (r *SpotifyTokenRepo) FindByInternalUserID(ctx context.Context, internalID string) (*model.SpotifyToken, error) {
  const q = `
    SELECT spotify_user_id, access_token, refresh_token, generated_at, expires_at
    FROM spotify_tokens
    WHERE internal_user_id = ?
  `
  t := &model.SpotifyToken{InternalUserID: internalID}
  row := r.db.QueryRowContext(ctx, q, internalID)
  if err := row.Scan(&t.SpotifyUserID, &t.AccessToken, &t.RefreshToken, &t.GeneratedAt, &t.ExpiresAt); err != nil {
    return nil, err
  }
  return t, nil
}
