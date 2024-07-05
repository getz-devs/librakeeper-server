package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/getz-devs/librakeeper-server/internal/server/models"
	"github.com/getz-devs/librakeeper-server/internal/server/services/books"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
)

type BookHandlers struct {
	service *books.BookService
	log     *slog.Logger
}

func NewBookHandlers(service *books.BookService, log *slog.Logger) *BookHandlers {
	return &BookHandlers{
		service: service,
		log:     log,
	}
}

func (h *BookHandlers) CreateBook(c *gin.Context) {
	var book models.Book
	if err := c.BindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ctx := context.WithValue(c.Request.Context(), "userID", userID)

	if err := h.service.Create(ctx, &book); err != nil {
		h.log.Error("failed to create book", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
		return
	}

	c.JSON(http.StatusCreated, book)
}

func (h *BookHandlers) GetBook(c *gin.Context) {
	bookID := c.Param("id")

	ctx := c.Request.Context()
	book, err := h.service.GetByID(ctx, bookID)
	if err != nil {
		if errors.Is(err, books.ErrBookNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.log.Error("failed to get book", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get book"})
		return
	}

	c.JSON(http.StatusOK, book)
}

func (h *BookHandlers) GetBooksByUserID(c *gin.Context) {
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
		h.log.Error("failed to get result", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get result"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *BookHandlers) GetBooksByBookshelfID(c *gin.Context) {
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
		h.log.Error("failed to get result by bookshelf id", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get result by bookshelf ID"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *BookHandlers) UpdateBook(c *gin.Context) {
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
		if errors.Is(err, books.ErrBookNotFound) {
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

func (h *BookHandlers) DeleteBook(c *gin.Context) {
	bookID := c.Param("id")

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ctx := context.WithValue(c.Request.Context(), "userID", userID)

	if err := h.service.Delete(ctx, bookID); err != nil {
		if errors.Is(err, books.ErrBookNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.log.Error("failed to delete book", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}
