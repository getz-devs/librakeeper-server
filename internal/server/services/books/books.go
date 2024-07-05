package books

import (
	"context"
	"fmt"
	"github.com/getz-devs/librakeeper-server/internal/server/models"
	"github.com/getz-devs/librakeeper-server/internal/server/repository"
	"log/slog"
)

// BookService defines the interface for book service operations.
type BookService struct {
	repo          repository.BookRepo
	bookshelfRepo repository.BookshelfRepo
	log           *slog.Logger
	bookLimit     int // Limit for books per bookshelf
}

// NewBookService creates a new BookService instance.
func NewBookService(repo repository.BookRepo, bookshelfRepo repository.BookshelfRepo, log *slog.Logger) *BookService {
	return &BookService{
		repo:          repo,
		bookshelfRepo: bookshelfRepo,
		log:           log,
		bookLimit:     1000, // Hardcoded for now, TODO: read from config
	}
}

// Create creates a new book.
func (s *BookService) Create(ctx context.Context, book *models.Book) error {
	// Rule 2: Book Title & Author Presence
	if book.Title == "" || book.Author == "" {
		return fmt.Errorf("book title and author are required")
	}

	// Rule 3: Bookshelf Ownership
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return fmt.Errorf("userID not found in context")
	}

	bookshelf, err := s.bookshelfRepo.GetByID(ctx, book.BookshelfID)
	if err != nil {
		return fmt.Errorf("failed to get bookshelf: %w", err)
	}
	if bookshelf.UserID != userID {
		return fmt.Errorf("user is not authorized to add a book to this bookshelf")
	}

	// TODO: optimize
	// Rule 4: Book Limit per Bookshelf
	books, err := s.repo.GetByBookshelfID(ctx, book.BookshelfID, 1, int64(s.bookLimit)) // Get up to the limit
	if err != nil {
		return fmt.Errorf("failed to get books for bookshelf: %w", err)
	}
	if len(books) >= s.bookLimit {
		return fmt.Errorf("bookshelf has reached the book limit (%d)", s.bookLimit)
	}

	// Rule 5: Unique Book within Bookshelf
	for _, existingBook := range books {
		if existingBook.ISBN == book.ISBN {
			return fmt.Errorf("book with ISBN '%s' already exists in this bookshelf", book.ISBN)
		}
	}

	if err := s.repo.Create(ctx, book); err != nil {
		return fmt.Errorf("failed to create book: %w", err)
	}
	return nil
}

// GetByID retrieves a book by its ID.
func (s *BookService) GetByID(ctx context.Context, bookID string) (*models.Book, error) {
	book, err := s.repo.GetByID(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get book: %w", err)
	}
	return book, nil
}

// GetByUserID retrieves a list of books for a specific user.
func (s *BookService) GetByUserID(ctx context.Context, userID string, page int64, limit int64) ([]*models.Book, error) {
	books, err := s.repo.GetByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get books by user ID: %w", err)
	}
	return books, nil
}

// GetByBookshelfID retrieves books by bookshelf ID.
func (s *BookService) GetByBookshelfID(ctx context.Context, bookshelfID string, page int64, limit int64) ([]*models.Book, error) {
	books, err := s.repo.GetByBookshelfID(ctx, bookshelfID, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get books by bookshelf ID: %w", err)
	}
	return books, nil
}

// Update updates an existing book.
func (s *BookService) Update(ctx context.Context, bookID string, update *models.BookUpdate) error {
	// 1. Get the book
	book, err := s.repo.GetByID(ctx, bookID)
	if err != nil {
		return fmt.Errorf("failed to get book: %w", err)
	}

	// 2. Get the bookshelf
	bookshelf, err := s.bookshelfRepo.GetByID(ctx, book.BookshelfID)
	if err != nil {
		return fmt.Errorf("failed to get bookshelf: %w", err)
	}

	// 3. Get userID from context
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return fmt.Errorf("userID not found in context")
	}

	// 4. Check bookshelf ownership
	if bookshelf.UserID != userID {
		return fmt.Errorf("user is not authorized to modify this book")
	}

	// 5. If authorized, proceed with the update:
	if err := s.repo.Update(ctx, bookID, update); err != nil {
		return fmt.Errorf("failed to update book: %w", err)
	}

	return nil
}

// Delete deletes a book.
func (s *BookService) Delete(ctx context.Context, bookID string) error {
	if err := s.repo.Delete(ctx, bookID); err != nil {
		return fmt.Errorf("failed to delete book: %w", err)
	}
	return nil
}
