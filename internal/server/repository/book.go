package repository

import (
	"context"
	"github.com/getz-devs/librakeeper-server/internal/server/models"
)

// BookRepo defines the interface for book repository operations.
type BookRepo interface {
	Create(ctx context.Context, book *models.Book) error
	GetByID(ctx context.Context, id string) (*models.Book, error)
	GetByISBN(ctx context.Context, isbn string) (*models.Book, error)
	GetByUserID(ctx context.Context, userID string, page int64, limit int64) ([]*models.Book, error)
	GetByBookshelfID(ctx context.Context, bookshelfID string, page int64, limit int64) ([]*models.Book, error)
	CountInBookshelf(ctx context.Context, bookshelfID string) (int, error)
	ExistsInBookshelf(ctx context.Context, isbn, bookshelfID string) (bool, error)
	Update(ctx context.Context, id string, update *models.BookUpdate) error
	Delete(ctx context.Context, id string) error
}
