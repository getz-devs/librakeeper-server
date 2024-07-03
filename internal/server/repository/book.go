package repository

import (
	"context"
	"github.com/getz-devs/librakeeper-server/internal/server/models"
)

// BookRepo defines the interface for book repository operations.
type BookRepo interface {
	Create(ctx context.Context, book *models.Book) error
	GetByID(ctx context.Context, id string) (*models.Book, error)
	GetByUserID(ctx context.Context, id string, page int64, limit int64) ([]*models.Book, error)
	GetByBookshelfID(ctx context.Context, id string, page int64, limit int64) ([]*models.Book, error)
	Update(ctx context.Context, id string, update *models.Book) error
	Delete(ctx context.Context, id string) error
}
