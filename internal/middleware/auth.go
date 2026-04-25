package middleware

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luponetn/lexpay/internal/utils"
)

// AuthMiddleware protects routes by requiring a valid access token in the cookies.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("access_token")
		if err != nil {
			slog.Error("access token not found", "error", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "failed",
				"message": "unauthorized: access token not found",
			})
			return
		}

		claims, err := utils.VerifyAccessToken(tokenString)
		if err != nil {
			slog.Error("invalid or expired access token", "error", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "failed",
				"message": "unauthorized: invalid or expired access token",
			})
			return
		}

		// Check if it's specifically an access token
		tokenType, _ := claims["type"].(string)
		if tokenType != "access" {
			slog.Warn("unauthorized: invalid token type", "type", tokenType)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "failed",
				"message": "unauthorized: invalid token type",
			})
			return
		}

		userID, ok := claims["userid"]
		if !ok {
			slog.Warn("unauthorized: user id not found in token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "failed",
				"message": "unauthorized: user id not found in token",
			})
			return
		}

		email, _ := claims["email"].(string)

		// Set user info in context for downstream handlers
		c.Set("userID", userID)
		c.Set("email", email)

		c.Next()
	}
}

// GetUserFromContext is a helper to retrieve user info in handlers
func GetUserFromContext(c *gin.Context) (interface{}, string, error) {
	userID, existsID := c.Get("userID")
	email, existsEmail := c.Get("email")

	if !existsID || !existsEmail {
		return nil, "", errors.New("user context not found")
	}

	return userID, email.(string), nil
}
