package middlewares

import (
	"context"
	"github.com/getz-devs/librakeeper-server/internal/server/services/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a Gin middleware for authenticating requests using Firebase
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Strip the "Bearer " prefix from the Authorization header value
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}
		idToken := strings.TrimPrefix(authHeader, "Bearer ")

		// Verify the ID token
		token, err := auth.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": err.Error()})
			c.Abort()
			return
		}

		// Add the UID to the context for further use.
		ctx := context.WithValue(c.Request.Context(), "userID", token.UID)
		c.Request = c.Request.WithContext(ctx) // Update the request's context

		c.Next()
	}
}
