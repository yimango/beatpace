package services

import (
	"encoding/base64"
	"log"
	"os"
	"time"

	_ "github.com/yimango/beatpace-backend/envloader"

	"github.com/golang-jwt/jwt/v5"
)

// jwtSecret holds the decoded JWT signing key.
var jwtSecret []byte

func init() {
	raw := os.Getenv("JWT_SECRET")
	if raw == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}
	// The secret is stored base64-encoded in the environment, so decode it
	decoded, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		log.Fatalf("failed to base64-decode JWT_SECRET: %v", err)
	}
	jwtSecret = decoded
	log.Printf("[AUTH] JWT secret loaded, %d bytes", len(jwtSecret))
}

// GenerateJWTToken creates an HS256-signed JWT with `sub` set to the given userID.
func GenerateJWTToken(userID string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
