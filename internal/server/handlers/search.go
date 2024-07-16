package handlers

import (
	"errors"
	"github.com/getz-devs/librakeeper-server/internal/server/services/search"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

// SearchHandlers handles HTTP requests related to search.
type SearchHandlers struct {
	service *search.SearchService
	log     *slog.Logger
}

func NewSearchHandlers(service *search.SearchService, log *slog.Logger) *SearchHandlers {
	return &SearchHandlers{
		service: service,
		log:     log,
	}
}

func (s *SearchHandlers) Simple(c *gin.Context) {
	isbn := c.Query("isbn")
	ctx := c.Request.Context()
	resp, err := s.service.Simple(ctx, isbn)
	if err != nil {
		if errors.Is(err, search.ErrISBNNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		s.log.Error("failed to get bookshelf", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get bookshelf"})
		return
	}

	c.JSON(http.StatusOK, resp)

}

func (s *SearchHandlers) Advanced(c *gin.Context) {
	const op = "handlers.SearchHandlers.Advanced"
	log := s.log.With(slog.String("op", op))
	isbn := c.Query("isbn")
	ctx := c.Request.Context()
	resp, err := s.service.Advanced(ctx, isbn)
	if err != nil {
		if errors.Is(err, search.ErrISBNNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		log.Error("failed to search", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search"})
		return
	}
	c.JSON(http.StatusOK, resp)
}
