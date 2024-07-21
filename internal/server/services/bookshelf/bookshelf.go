package bookshelf

import (
	"context"
	"errors"
	"fmt"
	"github.com/getz-devs/librakeeper-server/internal/server/models"
	"github.com/getz-devs/librakeeper-server/internal/server/repository"
	"github.com/getz-devs/librakeeper-server/internal/server/storage/mongo"
	"log/slog"
)

// Custom Error Types:
var (
	ErrBookshelfNotFound      = errors.New("bookshelf not found")
	ErrNameRequired           = errors.New("bookshelf name is required")
	ErrUserNotFoundInContext  = errors.New("userID not found in context")
	ErrNotAuthorized          = errors.New("user is not authorized to perform this action")
	ErrBookshelfAlreadyExists = errors.New("bookshelf with this name already exists for this user")
)

// BookshelfService handles business logic for bookshelf.
type BookshelfService struct {
	repo repository.BookshelfRepo
	log  *slog.Logger
}

// NewBookshelfService creates a new BookshelfService instance.
func NewBookshelfService(repo repository.BookshelfRepo, log *slog.Logger) *BookshelfService {
	return &BookshelfService{
		repo: repo,
		log:  log,
	}
}

// Create a new bookshelf.
func (s *BookshelfService) Create(ctx context.Context, bookshelf *models.Bookshelf) error {
	// Rule 1: Bookshelf Name Presence
	if bookshelf.Name == "" {
		return ErrNameRequired
	}

	// Get userID from context
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return ErrUserNotFoundInContext
	}

	// Rule 2: Unique Bookshelf Name per User
	exists, err := s.repo.ExistsByNameAndUser(ctx, bookshelf.Name, userID)
	if err != nil {
		return fmt.Errorf("failed to check bookshelf existence: %w", err)
	}
	if exists {
		return ErrBookshelfAlreadyExists
	}

	// Set the UserID for the bookshelf
	bookshelf.UserID = userID

	if err := s.repo.Create(ctx, bookshelf); err != nil {
		return fmt.Errorf("failed to create bookshelf: %w", err)
	}

	return nil
}

// GetByID retrieves a bookshelf by its ID.
func (s *BookshelfService) GetByID(ctx context.Context, bookshelfID string) (*models.Bookshelf, error) {
	bookshelf, err := s.repo.GetByID(ctx, bookshelfID)
	if err != nil {
		if errors.Is(err, mongo.ErrBookshelfNotFound) {
			return nil, ErrBookshelfNotFound
		}
		return nil, fmt.Errorf("failed to get bookshelf: %w", err)
	}
	return bookshelf, nil
}

// GetByUser retrieves a list of bookshelf for a specific user.
func (s *BookshelfService) GetByUser(ctx context.Context, userID string, page int64, limit int64) ([]*models.Bookshelf, error) {
	bookshelves, err := s.repo.GetByUser(ctx, userID, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get bookshelf by user ID: %w", err)
	}
	return bookshelves, nil
}

// Update updates an existing bookshelf.
func (s *BookshelfService) Update(ctx context.Context, bookshelfID string, update *models.BookshelfUpdate) error {
	// Get the bookshelf
	bookshelf, err := s.repo.GetByID(ctx, bookshelfID)
	if err != nil {
		if errors.Is(err, mongo.ErrBookshelfNotFound) {
			return ErrBookshelfNotFound
		}
		return fmt.Errorf("failed to get bookshelf: %w", err)
	}

	// Get userID from context
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return ErrUserNotFoundInContext
	}

	// Check bookshelf ownership
	if bookshelf.UserID != userID {
		return ErrNotAuthorized
	}

	if err := s.repo.Update(ctx, bookshelfID, update); err != nil {
		return fmt.Errorf("failed to update bookshelf: %w", err)
	}

	return nil
}

// Delete deletes a bookshelf.
func (s *BookshelfService) Delete(ctx context.Context, bookshelfID string) error {
	// Get the bookshelf
	bookshelf, err := s.repo.GetByID(ctx, bookshelfID)
	if err != nil {
		if errors.Is(err, mongo.ErrBookshelfNotFound) {
			return ErrBookshelfNotFound
		}
		return fmt.Errorf("failed to get bookshelf: %w", err)
	}

	// Get userID from context
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return ErrUserNotFoundInContext
	}

	// Check bookshelf ownership
	if bookshelf.UserID != userID {
		return ErrNotAuthorized
	}

	if err := s.repo.Delete(ctx, bookshelfID); err != nil {
		return fmt.Errorf("failed to delete bookshelf: %w", err)
	}

	return nil
}
