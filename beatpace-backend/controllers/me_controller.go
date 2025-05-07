package controllers
import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func MeHandler(c *gin.Context) {
	uid, exists := c.Get("userID")
	if !exists {
	  c.JSON(http.StatusUnauthorized, gin.H{"error":"not logged in"})
	  return
	}
	// Optional: fetch spotifyID or any other info from your DB
	c.JSON(http.StatusOK, gin.H{"userID": uid})
  }
  