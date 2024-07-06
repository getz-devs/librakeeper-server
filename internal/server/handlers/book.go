package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/getz-devs/librakeeper-server/internal/server/models"
	"github.com/getz-devs/librakeeper-server/internal/server/services/book"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
)

// BookHandlers handles HTTP requests related to books.
type BookHandlers struct {
	service *book.BookService
	log     *slog.Logger
}

// NewBookHandlers creates a new BookHandlers instance.
func NewBookHandlers(service *book.BookService, log *slog.Logger) *BookHandlers {
	return &BookHandlers{
		service: service,
		log:     log,
	}
}

// Create handles the creation of a new book.
func (h *BookHandlers) Create(c *gin.Context) {
	var b models.Book
	if err := c.BindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ctx := context.WithValue(c.Request.Context(), "userID", userID)

	if err := h.service.Create(ctx, &b); err != nil {
		h.log.Error("failed to create book", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
		return
	}

	c.JSON(http.StatusCreated, b)
}

// GetByID retrieves a book by ID.
func (h *BookHandlers) GetByID(c *gin.Context) {
	bookID := c.Param("id")

	ctx := c.Request.Context()
	b, err := h.service.GetByID(ctx, bookID)
	if err != nil {
		if errors.Is(err, book.ErrBookNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.log.Error("failed to get book", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get book"})
		return
	}

	c.JSON(http.StatusOK, b)
}

// GetByUser retrieves books for a user.
func (h *BookHandlers) GetByUser(c *gin.Context) {
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

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ctx := c.Request.Context()
	result, err := h.service.GetByUserID(ctx, userID.(string), page, limit)
	if err != nil {
		h.log.Error("failed to get books by user ID", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get books by user ID"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetByBookshelfID retrieves books from a specific bookshelf.
func (h *BookHandlers) GetByBookshelfID(c *gin.Context) {
	bookshelfID := c.Param("bookshelfId")

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

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ctx := context.WithValue(c.Request.Context(), "userID", userID)

	result, err := h.service.GetByBookshelfID(ctx, bookshelfID, page, limit)
	if err != nil {
		h.log.Error("failed to get books by bookshelf ID", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get books by bookshelf ID"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Update handles updating a book.
func (h *BookHandlers) Update(c *gin.Context) {
	bookID := c.Param("id")

	var update models.BookUpdate
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

	if err := h.service.Update(ctx, bookID, &update); err != nil {
		if errors.Is(err, book.ErrBookNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		h.log.Error(
			"failed to update book",
			slog.Any("error", err),
			slog.String("bookID", bookID),
			slog.String("userID", fmt.Sprintf("%v", userID)),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book updated successfully"})
}

// Delete handles the deletion of a book.
func (h *BookHandlers) Delete(c *gin.Context) {
	bookID := c.Param("id")

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ctx := context.WithValue(c.Request.Context(), "userID", userID)

	if err := h.service.Delete(ctx, bookID); err != nil {
		if errors.Is(err, book.ErrBookNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.log.Error("failed to delete book", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}
