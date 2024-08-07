package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/getz-devs/librakeeper-server/internal/server/models"
	"github.com/getz-devs/librakeeper-server/internal/server/services/bookshelf"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
)

// BookshelfHandlers handles HTTP requests related to bookshelf.
type BookshelfHandlers struct {
	service *bookshelf.BookshelfService
	log     *slog.Logger
}

// NewBookshelfHandlers creates a new BookshelfHandlers instance.
func NewBookshelfHandlers(service *bookshelf.BookshelfService, log *slog.Logger) *BookshelfHandlers {
	return &BookshelfHandlers{
		service: service,
		log:     log,
	}
}

// Create creates a new bookshelf.
func (h *BookshelfHandlers) Create(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var b models.Bookshelf
	if err := c.BindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.WithValue(c.Request.Context(), "userID", userID)
	if err := h.service.Create(ctx, &b); err != nil {
		if errors.Is(err, bookshelf.ErrNameRequired) || errors.Is(err, bookshelf.ErrBookshelfAlreadyExists) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		h.log.Error("failed to create bookshelf", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create bookshelf"})
		return
	}

	c.JSON(http.StatusCreated, b)
}

// GetByID retrieves a bookshelf by ID.
func (h *BookshelfHandlers) GetByID(c *gin.Context) {
	bookshelfID := c.Param("id")

	ctx := c.Request.Context()
	b, err := h.service.GetByID(ctx, bookshelfID)
	if err != nil {
		if errors.Is(err, bookshelf.ErrBookshelfNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.log.Error("failed to get bookshelf", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get bookshelf"})
		return
	}

	c.JSON(http.StatusOK, b)
}

// GetByUser retrieves bookshelf for a user.
func (h *BookshelfHandlers) GetByUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
		return
	}

	ctx := c.Request.Context()
	result, err := h.service.GetByUser(ctx, userID.(string), page, limit)
	if err != nil {
		h.log.Error("failed to get result by user id", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get bookshelf by user ID"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Update updates a bookshelf.
func (h *BookshelfHandlers) Update(c *gin.Context) {
	bookshelfID := c.Param("id")

	var update models.BookshelfUpdate
	if err := c.BindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ctx := context.WithValue(c.Request.Context(), "userID", userID)

	if err := h.service.Update(ctx, bookshelfID, &update); err != nil {
		if errors.Is(err, bookshelf.ErrBookshelfNotFound) || errors.Is(err, bookshelf.ErrNotAuthorized) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		h.log.Error(
			"failed to update bookshelf",
			slog.Any("error", err),
			slog.String("bookshelfID", bookshelfID),
			slog.String("userID", fmt.Sprintf("%v", userID)),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update bookshelf"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bookshelf updated successfully"})
}

// Delete deletes a bookshelf.
func (h *BookshelfHandlers) Delete(c *gin.Context) {
	bookshelfID := c.Param("id")

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ctx := context.WithValue(c.Request.Context(), "userID", userID)

	if err := h.service.Delete(ctx, bookshelfID); err != nil {
		if errors.Is(err, bookshelf.ErrBookshelfNotFound) || errors.Is(err, bookshelf.ErrNotAuthorized) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.log.Error("failed to delete bookshelf", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete bookshelf"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bookshelf deleted successfully"})
}
