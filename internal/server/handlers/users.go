package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/getz-devs/librakeeper-server/internal/server/models"
	"github.com/getz-devs/librakeeper-server/internal/server/services/users"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

// UserHandlers handles HTTP requests related to users.
type UserHandlers struct {
	service *users.UserService
	log     *slog.Logger
}

// NewUserHandlers creates a new UserHandlers instance.
func NewUserHandlers(service *users.UserService, log *slog.Logger) *UserHandlers {
	return &UserHandlers{
		service: service,
		log:     log,
	}
}

// Create creates a new user.
func (h *UserHandlers) Create(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.ID = userID.(string) // Assign Firebase UID to user.ID

	ctx := c.Request.Context()

	if err := h.service.Create(ctx, &user); err != nil {
		if errors.Is(err, users.ErrNotAuthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		h.log.Error("failed to create user", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// GetByID retrieves a user by ID.
func (h *UserHandlers) GetByID(c *gin.Context) {
	userID := c.Param("id")

	ctx := c.Request.Context()
	user, err := h.service.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, users.ErrNotAuthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		h.log.Error("failed to get user", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Update updates a user's information.
func (h *UserHandlers) Update(c *gin.Context) {
	userID := c.Param("id")

	var update models.UserUpdate
	if err := c.BindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ctx := context.WithValue(c.Request.Context(), "userID", authUserID)

	if err := h.service.Update(ctx, userID, &update); err != nil {
		if errors.Is(err, users.ErrUserNotFound) || errors.Is(err, users.ErrNotAuthorized) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.log.Error(
			"failed to update user",
			slog.Any("error", err),
			slog.String("userID", userID),
			slog.String("authUserID", fmt.Sprintf("%v", authUserID)),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// Delete removes a user.
func (h *UserHandlers) Delete(c *gin.Context) {
	userID := c.Param("id")

	authUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ctx := context.WithValue(c.Request.Context(), "userID", authUserID)

	if err := h.service.Delete(ctx, userID); err != nil {
		if errors.Is(err, users.ErrUserNotFound) || errors.Is(err, users.ErrNotAuthorized) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.log.Error("failed to delete user", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
