package repository

import (
	"context"
	"github.com/getz-devs/librakeeper-server/internal/server/models"
)

// BookshelfRepo defines the interface for bookshelf repository operations.
type BookshelfRepo interface {
	Create(ctx context.Context, bookshelf *models.Bookshelf) error
	GetByID(ctx context.Context, id string) (*models.Bookshelf, error)
	GetByUser(ctx context.Context, userID string, page int64, limit int64) ([]*models.Bookshelf, error) // Renamed
	CountByUser(ctx context.Context, userID string) (int, error)
	ExistsByNameAndUser(ctx context.Context, name, userID string) (bool, error)
	Update(ctx context.Context, id string, update *models.BookshelfUpdate) error
	Delete(ctx context.Context, id string) error
}
