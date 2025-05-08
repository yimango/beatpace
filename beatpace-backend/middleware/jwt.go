package middleware

import (
  "net/http"
  "time"

  "github.com/gin-gonic/gin"
  "github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(/* your secret from env */)

func JWT() gin.HandlerFunc {
  return func(c *gin.Context) {
    tokenString, err := c.Cookie("app_jwt")
    if err != nil {
      c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error":"missing token"})
      return
    }

    tok, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
      return jwtSecret, nil
    })
    if err != nil || !tok.Valid {
      c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error":"invalid token"})
      return
    }

    claims := tok.Claims.(*jwt.RegisteredClaims)
    if claims.ExpiresAt.Time.Before(time.Now()) {
      c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error":"token expired"})
      return
    }

    c.Set("userID", claims.Subject)
    c.Next()
  }
}
