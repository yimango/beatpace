package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yimango/beatpace-backend/services"
)

// JWT middleware validates the session token and sets the user ID in the context
func JWT(authService services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("JWT middleware: checking authorization")
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			fmt.Println("JWT middleware: missing authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			fmt.Printf("JWT middleware: invalid authorization header format: %s\n", authHeader)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			c.Abort()
			return
		}

		token := parts[1]
		fmt.Printf("JWT middleware: validating token: %s\n", token)
		session, err := authService.ValidateSession(c.Request.Context(), token)
		if err != nil {
			fmt.Printf("JWT middleware: token validation failed: %v\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		// Set user ID in context for later use
		fmt.Printf("JWT middleware: setting userID in context: %s\n", session.UserID.String())
		c.Set("userID", session.UserID.String())
		c.Next()
	}
}
