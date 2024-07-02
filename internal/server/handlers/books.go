package handlers

import (
	"errors"
	"github.com/getz-devs/librakeeper-server/internal/server/models"
	"github.com/getz-devs/librakeeper-server/internal/server/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"
	"net/http"
	"strconv"
)

type BookHandlers struct {
	service *services.BookService
	log     *slog.Logger
}

func NewBookHandlers(service *services.BookService, log *slog.Logger) *BookHandlers {
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

	ctx := c.Request.Context()
	createdBook, err := h.service.CreateBook(ctx, &book)
	if err != nil {
		h.log.Error("failed to create book", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
		return
	}

	c.JSON(http.StatusCreated, createdBook)
}

func (h *BookHandlers) GetBook(c *gin.Context) {
	bookIDHex := c.Param("id")

	bookID, err := primitive.ObjectIDFromHex(bookIDHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	ctx := c.Request.Context()
	book, err := h.service.GetBook(ctx, bookID)
	if err != nil {
		if errors.Is(err, services.ErrBookNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.log.Error("failed to get book", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get book"})
		return
	}

	c.JSON(http.StatusOK, book)
}

func (h *BookHandlers) GetBooks(c *gin.Context) {
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
	books, err := h.service.GetBooks(ctx, page, limit)
	if err != nil {
		h.log.Error("failed to get books", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get books"})
		return
	}

	c.JSON(http.StatusOK, books)
}

func (h *BookHandlers) GetBooksByBookshelfID(c *gin.Context) {
	bookshelfIDHex := c.Param("bookshelfId")

	bookshelfID, err := primitive.ObjectIDFromHex(bookshelfIDHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bookshelf ID"})
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
	books, err := h.service.GetBooksByBookshelfID(ctx, bookshelfID, page, limit)
	if err != nil {
		h.log.Error("failed to get books by bookshelf id", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get books by bookshelf ID"})
		return
	}

	c.JSON(http.StatusOK, books)
}

func (h *BookHandlers) UpdateBook(c *gin.Context) {
	bookIDHex := c.Param("id")

	bookID, err := primitive.ObjectIDFromHex(bookIDHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	var update models.Book
	if err := c.BindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	err = h.service.UpdateBook(ctx, bookID, update.ToMap())
	if err != nil {
		if errors.Is(err, services.ErrBookNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.log.Error("failed to update book", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book updated successfully"})
}

func (h *BookHandlers) DeleteBook(c *gin.Context) {
	bookIDHex := c.Param("id")

	bookID, err := primitive.ObjectIDFromHex(bookIDHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	ctx := c.Request.Context()
	err = h.service.DeleteBook(ctx, bookID)
	if err != nil {
		if errors.Is(err, services.ErrBookNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.log.Error("failed to delete book", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}
