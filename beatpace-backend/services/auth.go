package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yimango/beatpace-backend/model"
	"github.com/yimango/beatpace-backend/repository"
)

type authService struct {
	userRepo  repository.UserRepository
	tokenRepo repository.TokenRepository
}

func NewAuthService(userRepo repository.UserRepository, tokenRepo repository.TokenRepository) *authService {
	return &authService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
	}
}

func (s *authService) CreateSession(ctx context.Context, userID string) (*model.Session, error) {
	// Generate a random session token
	token, err := generateToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}

	// Parse userID to UUID
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %v", err)
	}

	// Create a new session
	session := &model.Session{
		ID:        uuid.New(),
		UserID:    uid,
		Token:     token,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour), // Sessions expire after 24 hours
	}

	// Save the session
	if err := s.tokenRepo.SaveSession(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to save session: %v", err)
	}

	return session, nil
}

func (s *authService) ValidateSession(ctx context.Context, token string) (*model.Session, error) {
	fmt.Printf("Validating session token: %s\n", token)
	session, err := s.tokenRepo.GetSessionByToken(ctx, token)
	if err != nil {
		fmt.Printf("Failed to get session by token: %v\n", err)
		return nil, fmt.Errorf("session not found: %v", err)
	}

	if time.Now().After(session.ExpiresAt) {
		fmt.Printf("Session expired at %v\n", session.ExpiresAt)
		return nil, fmt.Errorf("session expired")
	}

	fmt.Printf("Session valid for user %s\n", session.UserID.String())
	return session, nil
}

func (s *authService) RevokeSession(ctx context.Context, token string) error {
	return s.tokenRepo.DeleteSession(ctx, token)
}

// generateToken creates a cryptographically secure random token
func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
} 