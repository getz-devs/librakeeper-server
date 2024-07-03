package handlers

import (
	"errors"
	"github.com/getz-devs/librakeeper-server/internal/server/models"
	"github.com/getz-devs/librakeeper-server/internal/server/services/users"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type UserHandlers struct {
	service *users.UserService
	log     *slog.Logger
}

func NewUserHandlers(service *users.UserService, log *slog.Logger) *UserHandlers {
	return &UserHandlers{
		service: service,
		log:     log,
	}
}

func (h *UserHandlers) CreateUser(c *gin.Context) {
	userID, exists := c.Get("uid")
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
	createdUser, err := h.service.CreateUser(ctx, &user)
	if err != nil {
		if errors.Is(err, users.ErrUserAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		h.log.Error("failed to create user", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}

func (h *UserHandlers) GetUser(c *gin.Context) {
	userID := c.Param("id") // Get userID directly from the URL

	ctx := c.Request.Context()
	user, err := h.service.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.log.Error("failed to get user", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandlers) UpdateUser(c *gin.Context) {
	userID := c.Param("id")

	var update models.User
	if err := c.BindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	err := h.service.UpdateUser(ctx, userID, update.ToMap())
	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.log.Error("failed to update user", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (h *UserHandlers) DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	ctx := c.Request.Context()
	err := h.service.DeleteUser(ctx, userID)
	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.log.Error("failed to delete user", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
