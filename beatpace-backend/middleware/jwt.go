package middleware

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"encoding/base64"
)

var jwtSecret []byte

func init() {
	// Load and decode the JWT secret (base64-encoded) from environment
	raw := os.Getenv("JWT_SECRET")
	if raw == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}
	decoded, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		log.Fatalf("failed to decode JWT_SECRET: %v", err)
	}
	jwtSecret = decoded
	log.Printf("[JWT] Loaded JWT secret, %d bytes", len(jwtSecret))
}

// JWT returns a Gin middleware that validates the JWT in the "app_jwt" cookie.
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1) Extract the raw token from the cookie
		rawToken, err := c.Cookie("app_jwt")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		// DEBUG: log first part of the token
		log.Printf("[JWT] raw token: %s...", rawToken[:10])

		// 2) Parse and validate the token
		tok, err := jwt.ParseWithClaims(rawToken, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil {
			log.Printf("[JWT] parse error: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		if !tok.Valid {
			log.Println("[JWT] token is not valid")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// 3) Check expiry
		claims := tok.Claims.(*jwt.RegisteredClaims)
		if claims.ExpiresAt.Time.Before(time.Now()) {
			log.Printf("[JWT] token expired at %v", claims.ExpiresAt.Time)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
			return
		}

		// 4) Store the userID in context for downstream handlers
		c.Set("userID", claims.Subject)
		c.Next()
	}
}
