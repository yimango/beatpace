package controllers
import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/yimango/beatpace-backend/services"

)

func Callback(c *gin.Context) {
	// 1) Read code
	code := c.Query("code")
	if code == "" {
	  c.JSON(http.StatusBadRequest, gin.H{"error":"missing code"})
	  return
	}
  
	// 2) Exchange for tokens
	tok, err := services.GetSpotifyAccessToken(code)
	
	if err != nil {
	  c.JSON(http.StatusBadGateway, gin.H{"error":"token exchange failed"})
	  return
	}
	expiresAt := time.Now().Add(time.Duration(tok.ExpiresIn)*time.Second)
  
	// 3) Fetch Spotify user ID
	spUser, err := services.FetchSpotifyUser(tok.AccessToken)
	if err != nil {
	  c.JSON(http.StatusBadGateway, gin.H{"error":"failed to fetch profile"})
	  return
	}
  
	// 4) Find or create your user
	internalUserID, err := services.FindOrCreateUser(spUser.ID)
	if err != nil {
	  c.JSON(http.StatusInternalServerError, gin.H{"error":"db error"})
	  return
	}
  
	// 5) Save tokens
	if err := services.SaveSpotifyTokens(internalUserID,
		 tok.AccessToken, tok.RefreshToken, expiresAt); err != nil {
	  c.JSON(http.StatusInternalServerError, gin.H{"error":"failed to save tokens"})
	  return
	}
  
	// 6) Mint JWT
	jwtStr, err := services.GenerateJWTToken(internalUserID)
	if err != nil {
	  c.JSON(http.StatusInternalServerError, gin.H{"error":"failed to generate JWT"})
	  return
	}
  
	// 7) Set HttpOnly cookie
	c.SetCookie("app_jwt", jwtStr, 3600*24, "/", "", true, true)
  
	// 8) Redirect to your front end
	c.Redirect(http.StatusSeeOther, "http://localhost:3000")
  }
  