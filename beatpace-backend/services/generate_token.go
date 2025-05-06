// generate JWT Token
package services
import (
	"fmt"
	"time"
	"github.com/golang-jwt/jwt/v5"

)
// GenerateJWTToken generates a JWT token for the user
func GenerateJWTToken(userID string) (string, error) {
	// Define the token expiration time
	expirationTime := time.Now().Add(24 * time.Hour) // Token valid for 24 hours
	claims := &jwt.RegisteredClaims{
		Issuer:    userID,
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	// Create the token using the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign the token with a secret key
	secretKey := []byte("your_secret_key") // Replace with your actual secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}
	// Return the signed token string
	return tokenString, nil
}