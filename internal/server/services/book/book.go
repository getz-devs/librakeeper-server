package book

import (
	"context"
	"errors"
	"fmt"
	"github.com/getz-devs/librakeeper-server/internal/server/models"
	"github.com/getz-devs/librakeeper-server/internal/server/repository"
	"github.com/getz-devs/librakeeper-server/internal/server/services/search"
	"github.com/getz-devs/librakeeper-server/internal/server/storage/mongo"
	"log/slog"
	"time"
)

// Custom Error Types:
var (
	ErrBookNotFound           = errors.New("book not found")
	ErrBookshelfNotFound      = errors.New("bookshelf not found")
	ErrUserNotFoundInContext  = errors.New("userID not found in context")
	ErrNotAuthorized          = errors.New("user is not authorized to perform this action")
	ErrTitleAndAuthorRequired = errors.New("book title and author are required")
	ErrBookshelfLimitReached  = errors.New("bookshelf has reached the book limit")
	ErrBookAlreadyExists      = errors.New("book with this ISBN already exists in this bookshelf")
	ErrCantAddToAllBooks      = errors.New("error adding book to all books")
	ErrBookAlreadyExistsInALL = errors.New("book already exist in all books")
)

// BookService defines the interface for book service operations.
type BookService struct {
	repo          repository.BookRepo
	allBooksRepo  repository.BookRepo
	bookshelfRepo repository.BookshelfRepo
	searcher      *search.SearchService
	log           *slog.Logger
	bookLimit     int
}

// NewBookService creates a new BookService instance.
func NewBookService(repo repository.BookRepo, allBooksRepo repository.BookRepo, bookshelfRepo repository.BookshelfRepo, searcher *search.SearchService, log *slog.Logger) *BookService {
	return &BookService{
		repo:          repo,
		allBooksRepo:  allBooksRepo,
		bookshelfRepo: bookshelfRepo,
		searcher:      searcher,
		log:           log,
		bookLimit:     1000, // TODO: Read from config
	}
}

// Create creates a new book.
func (s *BookService) Create(ctx context.Context, book *models.Book) error {
	// Rule 2: Book Title & Author Presence
	if book.Title == "" || book.Author == "" {
		return ErrTitleAndAuthorRequired
	}

	// Rule 3: Bookshelf Ownership
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return ErrUserNotFoundInContext
	}

	if book.BookshelfID != "" {
		// if bookshelf specified, check bookshelf

		bookshelf, err := s.bookshelfRepo.GetByID(ctx, book.BookshelfID)
		if err != nil {
			if errors.Is(err, mongo.ErrBookshelfNotFound) {
				return ErrBookshelfNotFound
			}
			return fmt.Errorf("failed to get bookshelf: %w", err)
		}

		if bookshelf.UserID != userID {
			return ErrNotAuthorized
		}

		// Rule 4: Book Limit per Bookshelf
		bookCount, err := s.repo.CountInBookshelf(ctx, book.BookshelfID)
		if err != nil {
			return fmt.Errorf("failed to get book count for bookshelf: %w", err)
		}
		if bookCount >= s.bookLimit {
			return ErrBookshelfLimitReached
		}

		// Rule 5: Unique Book within Bookshelf
		exists, err := s.repo.ExistsInBookshelf(ctx, book.ISBN, book.BookshelfID)
		if err != nil {
			return fmt.Errorf("failed to check book existence: %w", err)
		}
		if exists {
			return ErrBookAlreadyExists
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
		if errors.Is(err, mongo.ErrBookNotFound) {
			return nil, ErrBookNotFound
		}
		return nil, fmt.Errorf("failed to get book: %w", err)
	}
	return book, nil
}

// GetByISBN retrieves a book by its ISBN.
func (s *BookService) GetByISBN(ctx context.Context, isbn string) (*models.Book, error) {
	book, err := s.repo.GetByISBN(ctx, isbn)
	if err != nil {
		if errors.Is(err, mongo.ErrBookNotFound) {
			return nil, ErrBookNotFound
		}
		return nil, fmt.Errorf("failed to get book by ISBN: %w", err)
	}
	return book, nil
}

// GetByUserID retrieves a list of book for a specific user.
func (s *BookService) GetByUserID(ctx context.Context, userID string, page int64, limit int64) ([]*models.Book, error) {
	books, err := s.repo.GetByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get book by user ID: %w", err)
	}
	return books, nil
}

// GetByBookshelfID retrieves book by bookshelf ID.
func (s *BookService) GetByBookshelfID(ctx context.Context, bookshelfID string, page int64, limit int64) ([]*models.Book, error) {
	books, err := s.repo.GetByBookshelfID(ctx, bookshelfID, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get book by bookshelf ID: %w", err)
	}
	return books, nil
}

// Update updates an existing book.
func (s *BookService) Update(ctx context.Context, bookID string, update *models.BookUpdate) error {
	// 1. Get the book
	book, err := s.repo.GetByID(ctx, bookID)
	if err != nil {
		if errors.Is(err, mongo.ErrBookNotFound) {
			return ErrBookNotFound
		}
		return fmt.Errorf("failed to get book: %w", err)
	}

	// 2. Get the bookshelf
	bookshelf, err := s.bookshelfRepo.GetByID(ctx, book.BookshelfID)
	if err != nil {
		if errors.Is(err, mongo.ErrBookshelfNotFound) {
			return ErrBookshelfNotFound
		}
		return fmt.Errorf("failed to get bookshelf: %w", err)
	}

	// 3. Get userID from context
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return ErrUserNotFoundInContext
	}

	// 4. Check bookshelf ownership
	if bookshelf.UserID != userID {
		return ErrNotAuthorized
	}

	// 5. If authorized, proceed with the update:
	if err := s.repo.Update(ctx, bookID, update); err != nil {
		return fmt.Errorf("failed to update book: %w", err)
	}

	return nil
}

// Delete deletes a book.
func (s *BookService) Delete(ctx context.Context, bookID string) error {
	// 1. Get the book
	book, err := s.repo.GetByID(ctx, bookID)
	if err != nil {
		if errors.Is(err, mongo.ErrBookNotFound) {
			return ErrBookNotFound
		}
		return fmt.Errorf("failed to get book: %w", err)
	}

	// 2. Get the bookshelf
	bookshelf, err := s.bookshelfRepo.GetByID(ctx, book.BookshelfID)
	if err != nil {
		if errors.Is(err, mongo.ErrBookshelfNotFound) {
			return ErrBookshelfNotFound
		}
		return fmt.Errorf("failed to get bookshelf: %w", err)
	}

	// 3. Get userID from context
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return ErrUserNotFoundInContext
	}

	// 4. Check bookshelf ownership
	if bookshelf.UserID != userID {
		return ErrNotAuthorized
	}

	// 5. If authorized, proceed to delete:
	if err := s.repo.Delete(ctx, bookID); err != nil {
		return fmt.Errorf("failed to delete book: %w", err)
	}

	return nil
}

func (s *BookService) AddAdvanced(ctx context.Context, isbn string, index int) error {
	resp, err := s.searcher.Advanced(ctx, isbn)
	if err != nil {
		if errors.Is(err, search.ErrISBNNotFound) {
			return fmt.Errorf("nothing found by isbn: %w", err)
		}
		s.log.Error("failed to search", slog.Any("error", err))
		return fmt.Errorf("failed to search: %w", err)
	}

	if index >= len(resp.Books) {
		return errors.New("invalid index")
	}

	// Добавляем книгу в allBooksRepo
	// Проверяем, существует ли книга по ISBN
	_, err = s.allBooksRepo.GetByISBN(ctx, isbn)
	if errors.Is(err, mongo.ErrBookNotFound) {
		// Книги нет в allBooksRepo, можно добавить
		book := resp.Books[index]
		allBook := &models.Book{
			Title:       book.Title,
			Author:      book.Author,
			ISBN:        book.ISBN,
			Publishing:  book.Publishing,
			Description: book.Description,
			CoverImage:  book.CoverImage,
			ShopName:    book.ShopName,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			// UserID и BookshelfID не устанавливаем
		}

		if err := s.allBooksRepo.Create(ctx, allBook); err != nil {
			s.log.Error("failed to create book in allBooksRepo", slog.Any("error", err))
			return ErrCantAddToAllBooks
		}

		return nil
	}
	if err != nil {
		s.log.Error("failed to check for existing book in allBooksRepo", slog.Any("error", err))
		return fmt.Errorf("failed to check for existing book in allBooksRepo: %w", err)
	}

	return ErrBookAlreadyExistsInALL
	// Если книга уже существует, ничего не делаем
}
