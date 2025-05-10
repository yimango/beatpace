// controllers/users_controllers.go
package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/yimango/beatpace-backend/services"
)

func Callback(c *gin.Context) {
    // 1) Read code
    code := c.Query("code")
    if code == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "missing code"})
        return
    }

    // 2) Exchange for tokens
    tok, err := services.GetSpotifyAccessToken(code)
    if err != nil {
        c.JSON(http.StatusBadGateway, gin.H{"error": "token exchange failed"})
        return
    }

    // 3) Fetch Spotify user ID
    fetchUserService := services.FetchUserService{AccessToken: tok.AccessToken}
    spotifyID, err := fetchUserService.FetchUserID()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user ID"})
        return
    }

    // 4) (optional) Save tokens in DB
    // if err := services.SaveSpotifyTokens(spotifyID, tok); err != nil {
    //     c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save tokens"})
    //     return
    // }

    // 5) Mint your JWT
    jwtToken, err := services.GenerateJWTToken(spotifyID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate JWT"})
        return
    }

    // 6) Set the cookie with SameSite=None
	c.SetCookie(
		"app_jwt",
		jwtToken,
		3600,           // e.g. 1 hour
		"/",            // path
		"", // domain
		true,           // secure
		true,           // httpOnly
	  )
	  // then redirect:
	  c.Redirect(http.StatusSeeOther, "http://localhost:3000/")
}
