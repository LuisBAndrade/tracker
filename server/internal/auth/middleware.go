package auth

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("session")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing session cookie"})
			return 
		}

		userID, err := ValidateJWT(tokenStr, os.Getenv("JWT_SECRET"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
			return 
		}

		c.Set("user_id", userID)
		c.Next()
	}
}