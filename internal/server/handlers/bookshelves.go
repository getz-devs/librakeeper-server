package handlers

import (
	"errors"
	"fmt"
	"github.com/getz-devs/librakeeper-server/internal/server/models"
	"github.com/getz-devs/librakeeper-server/internal/server/services/bookshelves"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"
	"net/http"
	"strconv"
)

type BookshelfHandlers struct {
	service *bookshelves.BookshelfService
	log     *slog.Logger
}

func NewBookshelfHandlers(service *bookshelves.BookshelfService, log *slog.Logger) *BookshelfHandlers {
	return &BookshelfHandlers{
		service: service,
		log:     log,
	}
}

func (h *BookshelfHandlers) CreateBookshelf(c *gin.Context) {
	userID, exists := c.Get("uid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var bookshelf models.Bookshelf
	if err := c.BindJSON(&bookshelf); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bookshelf.UserID = userID.(string) // Type assertion to string
	ctx := c.Request.Context()
	createdBookshelf, err := h.service.CreateBookshelf(ctx, &bookshelf)
	if err != nil {
		h.log.Error("failed to create bookshelf", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create bookshelf"})
		return
	}

	c.JSON(http.StatusCreated, createdBookshelf)
}

func (h *BookshelfHandlers) GetBookshelf(c *gin.Context) {
	bookshelfIDHex := c.Param("id")

	bookshelfID, err := primitive.ObjectIDFromHex(bookshelfIDHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bookshelf ID"})
		return
	}

	ctx := c.Request.Context()
	bookshelf, err := h.service.GetBookshelf(ctx, bookshelfID)
	if err != nil {
		if errors.Is(err, bookshelves.ErrBookshelfNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.log.Error("failed to get bookshelf", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get bookshelf"})
		return
	}

	c.JSON(http.StatusOK, bookshelf)
}

func (h *BookshelfHandlers) GetBookshelvesByUserID(c *gin.Context) {
	userIDHex, exists := c.Get("uid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", userIDHex))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing user id"})
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
	result, err := h.service.GetBookshelvesByUserID(ctx, userID, page, limit)
	if err != nil {
		h.log.Error("failed to get result by user id", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get result by user ID"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *BookshelfHandlers) UpdateBookshelf(c *gin.Context) {
	bookshelfIDHex := c.Param("id")

	bookshelfID, err := primitive.ObjectIDFromHex(bookshelfIDHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bookshelf ID"})
		return
	}

	var update models.Bookshelf
	if err := c.BindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	err = h.service.UpdateBookshelf(ctx, bookshelfID, update.ToMap())
	if err != nil {
		if errors.Is(err, bookshelves.ErrBookshelfNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.log.Error("failed to update bookshelf", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update bookshelf"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bookshelf updated successfully"})
}

func (h *BookshelfHandlers) DeleteBookshelf(c *gin.Context) {
	bookshelfIDHex := c.Param("id")

	bookshelfID, err := primitive.ObjectIDFromHex(bookshelfIDHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bookshelf ID"})
		return
	}

	ctx := c.Request.Context()
	err = h.service.DeleteBookshelf(ctx, bookshelfID)
	if err != nil {
		if errors.Is(err, bookshelves.ErrBookshelfNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.log.Error("failed to delete bookshelf", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete bookshelf"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bookshelf deleted successfully"})
}
