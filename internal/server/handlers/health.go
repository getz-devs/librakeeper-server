package handlers

import (
	"context"
	"github.com/getz-devs/librakeeper-server/internal/server/services/storage"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthCheck performs a more comprehensive health check.
func HealthCheck(c *gin.Context) {
	// Check Database Connectivity
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) // Short timeout for DB check
	defer cancel()

	err := storage.Ping(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "DOWN", "error": "Database connection failed"})
		return
	}

	// TODO: Add other checks here as needed (e.g., external services, dependencies)

	// If all checks pass
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}
