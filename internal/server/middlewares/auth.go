package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/getz-devs/librakeeper-server/internal/server/services"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
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
	token, err := services.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	// Set the UID in the context for further use
	c.Set("uid", token.UID)
	c.Next()
}