package book

import (
	"context"
	"github.com/getz-devs/librakeeper-server/internal/server/models"
	"github.com/getz-devs/librakeeper-server/internal/server/services/search"
	"github.com/getz-devs/librakeeper-server/internal/server/storage/mongo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"os"
	"testing"
	"time"
)

// MockRepository is a mock implementation of the repository.BookRepo interface.
type MockRepository struct {
	mock.Mock
}

// Create mocks the Create method of the BookRepo interface.
func (m *MockRepository) Create(ctx context.Context, book *models.Book) error {
	args := m.Called(ctx, book)
	return args.Error(0)
}

// GetByID mocks the GetByID method of the BookRepo interface.
func (m *MockRepository) GetByID(ctx context.Context, id string) (*models.Book, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Book), args.Error(1)
}

// GetByISBN mocks the GetByISBN method of the BookRepo interface.
func (m *MockRepository) GetByISBN(ctx context.Context, isbn string) (*models.Book, error) {
	args := m.Called(ctx, isbn)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Book), args.Error(1)
}

// GetByUserID mocks the GetByUserID method of the BookRepo interface.
func (m *MockRepository) GetByUserID(ctx context.Context, userID string, page int64, limit int64) ([]*models.Book, error) {
	args := m.Called(ctx, userID, page, limit)
	return args.Get(0).([]*models.Book), args.Error(1)
}

// GetByBookshelfID mocks the GetByBookshelfID method of the BookRepo interface.
func (m *MockRepository) GetByBookshelfID(ctx context.Context, bookshelfID string, page int64, limit int64) ([]*models.Book, error) {
	args := m.Called(ctx, bookshelfID, page, limit)
	return args.Get(0).([]*models.Book), args.Error(1)
}

// CountInBookshelf mocks the CountInBookshelf method of the BookRepo interface.
func (m *MockRepository) CountInBookshelf(ctx context.Context, bookshelfID string) (int, error) {
	args := m.Called(ctx, bookshelfID)
	return args.Int(0), args.Error(1)
}

// ExistsInBookshelf mocks the ExistsInBookshelf method of the BookRepo interface.
func (m *MockRepository) ExistsInBookshelf(ctx context.Context, isbn, bookshelfID string) (bool, error) {
	args := m.Called(ctx, isbn, bookshelfID)
	return args.Bool(0), args.Error(1)
}

// Update mocks the Update method of the BookRepo interface.
func (m *MockRepository) Update(ctx context.Context, id string, update *models.BookUpdate) error {
	args := m.Called(ctx, id, update)
	return args.Error(0)
}

// Delete mocks the Delete method of the BookRepo interface.
func (m *MockRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockBookshelfRepository is a mock implementation of the repository.BookshelfRepo interface.
type MockBookshelfRepository struct {
	mock.Mock
}

// Create mocks the Create method of the BookshelfRepo interface.
func (m *MockBookshelfRepository) Create(ctx context.Context, bookshelf *models.Bookshelf) error {
	args := m.Called(ctx, bookshelf)
	return args.Error(0)
}

// GetByID mocks the GetByID method of the BookshelfRepo interface.
func (m *MockBookshelfRepository) GetByID(ctx context.Context, id string) (*models.Bookshelf, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Bookshelf), args.Error(1)
}

// GetByUser mocks the GetByUser method of the BookshelfRepo interface.
func (m *MockBookshelfRepository) GetByUser(ctx context.Context, userID string, page int64, limit int64) ([]*models.Bookshelf, error) {
	args := m.Called(ctx, userID, page, limit)
	return args.Get(0).([]*models.Bookshelf), args.Error(1)
}

// CountByUser mocks the CountByUser method of the BookshelfRepo interface.
func (m *MockBookshelfRepository) CountByUser(ctx context.Context, userID string) (int, error) {
	args := m.Called(ctx, userID)
	return args.Int(0), args.Error(1)
}

// ExistsByNameAndUser mocks the ExistsByNameAndUser method of the BookshelfRepo interface.
func (m *MockBookshelfRepository) ExistsByNameAndUser(ctx context.Context, name, userID string) (bool, error) {
	args := m.Called(ctx, name, userID)
	return args.Bool(0), args.Error(1)
}

// Update mocks the Update method of the BookshelfRepo interface.
func (m *MockBookshelfRepository) Update(ctx context.Context, id string, update *models.BookshelfUpdate) error {
	args := m.Called(ctx, id, update)
	return args.Error(0)
}

// Delete mocks the Delete method of the BookshelfRepo interface.
func (m *MockBookshelfRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Helper function to create a pointer to a string
func stringPtr(s string) *string {
	return &s
}

// -- Tests -- //

func TestBookService_Create_Success(t *testing.T) {
	repo := new(MockRepository)
	bookshelfRepo := new(MockBookshelfRepository)
	searcher := new(search.SearchService)
	log := slog.New(slog.NewTextHandler(
		//io.Discard,
		//nil,
		os.Stdout,
		nil,
	))
	service := &BookService{
		repo:          repo,
		allBooksRepo:  repo, // In this test, both repos are the same mock
		bookshelfRepo: bookshelfRepo,
		searcher:      searcher,
		log:           log,
		bookLimit:     1000,
	}

	ctx := context.Background()
	userID := "testuser"
	ctx = context.WithValue(ctx, "userID", userID)

	book := &models.Book{
		UserID:      userID,
		BookshelfID: "testbookshelf",
		ISBN:        "1234567890",
		Title:       "Test Book",
		Author:      "Test Author",
	}

	// Mock the bookshelfRepo.GetByID to return a valid bookshelf
	bookshelfRepo.On("GetByID", ctx, book.BookshelfID).Return(
		&models.Bookshelf{
			ID:     book.BookshelfID,
			UserID: userID,
		}, nil,
	)

	// Mock the repo.CountInBookshelf to return a count less than the limit
	repo.On("CountInBookshelf", ctx, book.BookshelfID).Return(500, nil)

	// Mock the repo.ExistsInBookshelf to return false (book doesn't exist)
	repo.On("ExistsInBookshelf", ctx, book.ISBN, book.BookshelfID).Return(false, nil)

	// Mock the repo.Create to return no error
	repo.On("Create", ctx, book).Return(nil)

	err := service.Create(ctx, book)
	assert.NoError(t, err)

	repo.AssertExpectations(t)
	bookshelfRepo.AssertExpectations(t)
}

func TestBookService_Create_ErrorTitleAndAuthorRequired(t *testing.T) {
	repo := new(MockRepository)
	bookshelfRepo := new(MockBookshelfRepository)
	searcher := new(search.SearchService)
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	service := &BookService{
		repo:          repo,
		allBooksRepo:  repo,
		bookshelfRepo: bookshelfRepo,
		searcher:      searcher,
		log:           log,
		bookLimit:     1000,
	}

	// Test cases for missing title and/or author
	testCases := []struct {
		name  string
		book  *models.Book
		error error
	}{
		{
			name: "Missing Title",
			book: &models.Book{
				UserID:      "testuser",
				BookshelfID: "testbookshelf",
				ISBN:        "1234567890",
				Author:      "Test Author",
			},
			error: ErrTitleAndAuthorRequired,
		},
		{
			name: "Missing Author",
			book: &models.Book{
				UserID:      "testuser",
				BookshelfID: "testbookshelf",
				ISBN:        "1234567890",
				Title:       "Test Book",
			},
			error: ErrTitleAndAuthorRequired,
		},
		{
			name: "Missing Title and Author",
			book: &models.Book{
				UserID:      "testuser",
				BookshelfID: "testbookshelf",
				ISBN:        "1234567890",
			},
			error: ErrTitleAndAuthorRequired,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			userID := "testuser"
			ctx = context.WithValue(ctx, "userID", userID)

			err := service.Create(ctx, tc.book)
			assert.ErrorIs(t, err, tc.error)
		})
	}

	repo.AssertExpectations(t)
	bookshelfRepo.AssertExpectations(t)
}

func TestBookService_Create_ErrorUserNotFoundInContext(t *testing.T) {
	repo := new(MockRepository)
	bookshelfRepo := new(MockBookshelfRepository)
	searcher := new(search.SearchService)
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	service := &BookService{
		repo:          repo,
		allBooksRepo:  repo,
		bookshelfRepo: bookshelfRepo,
		searcher:      searcher,
		log:           log,
		bookLimit:     1000,
	}

	ctx := context.Background()

	book := &models.Book{
		UserID:      "testuser",
		BookshelfID: "testbookshelf",
		ISBN:        "1234567890",
		Title:       "Test Book",
		Author:      "Test Author",
	}

	err := service.Create(ctx, book)
	assert.ErrorIs(t, err, ErrUserNotFoundInContext)

	repo.AssertExpectations(t)
	bookshelfRepo.AssertExpectations(t)
}

func TestBookService_GetByID_Success(t *testing.T) {
	repo := new(MockRepository)
	bookshelfRepo := new(MockBookshelfRepository)
	searcher := new(search.SearchService)
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	service := &BookService{
		repo:          repo,
		allBooksRepo:  repo, // В этом тесте, allBooksRepo - тот же мок
		bookshelfRepo: bookshelfRepo,
		searcher:      searcher,
		log:           log,
		bookLimit:     1000,
	}

	ctx := context.Background()
	bookID := "testbookid"

	expectedBook := &models.Book{
		ID:          bookID,
		UserID:      "testuser",
		BookshelfID: "testbookshelf",
		ISBN:        "1234567890",
		Title:       "Test Book",
		Author:      "Test Author",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	repo.On("GetByID", ctx, bookID).Return(expectedBook, nil)

	book, err := service.GetByID(ctx, bookID)

	assert.NoError(t, err)
	assert.Equal(t, expectedBook, book)
	repo.AssertExpectations(t)
}

func TestBookService_GetByID_ErrorBookNotFound(t *testing.T) {
	repo := new(MockRepository)
	bookshelfRepo := new(MockBookshelfRepository)
	searcher := new(search.SearchService)
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	service := &BookService{
		repo:          repo,
		allBooksRepo:  repo,
		bookshelfRepo: bookshelfRepo,
		searcher:      searcher,
		log:           log,
		bookLimit:     1000,
	}

	ctx := context.Background()
	bookID := "nonexistentbookid"

	repo.On("GetByID", ctx, bookID).Return(nil, mongo.ErrBookNotFound)

	book, err := service.GetByID(ctx, bookID)

	assert.ErrorIs(t, err, ErrBookNotFound)
	assert.Nil(t, book)
	repo.AssertExpectations(t)
}

func TestBookService_Update_Success(t *testing.T) {
	repo := new(MockRepository)
	bookshelfRepo := new(MockBookshelfRepository)
	searcher := new(search.SearchService)
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	service := &BookService{
		repo:          repo,
		allBooksRepo:  repo,
		bookshelfRepo: bookshelfRepo,
		searcher:      searcher,
		log:           log,
		bookLimit:     1000,
	}

	ctx := context.Background()
	userID := "testuser"
	ctx = context.WithValue(ctx, "userID", userID)
	bookID := "testbookid"
	update := &models.BookUpdate{
		Title:       stringPtr("Updated Title"),
		Author:      stringPtr("Updated Author"),
		Publishing:  stringPtr("Updated Publishing"),
		Description: stringPtr("Updated Description"),
		CoverImage:  stringPtr("Updated CoverImage"),
		ShopName:    stringPtr("Updated ShopName"),
	}

	existingBook := &models.Book{
		ID:          bookID,
		UserID:      userID,
		BookshelfID: "testbookshelf",
		ISBN:        "1234567890",
		Title:       "Test Book",
		Author:      "Test Author",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	bookshelfRepo.On("GetByID", ctx, existingBook.BookshelfID).Return(
		&models.Bookshelf{
			ID:     existingBook.BookshelfID,
			UserID: userID,
		}, nil,
	)

	repo.On("GetByID", ctx, bookID).Return(existingBook, nil)
	repo.On("Update", ctx, bookID, update).Return(nil)

	err := service.Update(ctx, bookID, update)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
	bookshelfRepo.AssertExpectations(t)
}

func TestBookService_Update_ErrorBookNotFound(t *testing.T) {
	repo := new(MockRepository)
	bookshelfRepo := new(MockBookshelfRepository)
	searcher := new(search.SearchService)
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	service := &BookService{
		repo:          repo,
		allBooksRepo:  repo,
		bookshelfRepo: bookshelfRepo,
		searcher:      searcher,
		log:           log,
		bookLimit:     1000,
	}

	ctx := context.Background()
	userID := "testuser"
	ctx = context.WithValue(ctx, "userID", userID)
	bookID := "nonexistentbookid"
	update := &models.BookUpdate{
		Title: stringPtr("Updated Title"),
	}

	repo.On("GetByID", ctx, bookID).Return(nil, mongo.ErrBookNotFound)

	err := service.Update(ctx, bookID, update)

	assert.ErrorIs(t, err, ErrBookNotFound)
	repo.AssertExpectations(t)
}

func TestBookService_Delete_Success(t *testing.T) {
	repo := new(MockRepository)
	bookshelfRepo := new(MockBookshelfRepository)
	searcher := new(search.SearchService)
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	service := &BookService{
		repo:          repo,
		allBooksRepo:  repo,
		bookshelfRepo: bookshelfRepo,
		searcher:      searcher,
		log:           log,
		bookLimit:     1000,
	}

	ctx := context.Background()
	userID := "testuser"
	ctx = context.WithValue(ctx, "userID", userID)
	bookID := "testbookid"

	existingBook := &models.Book{
		ID:          bookID,
		UserID:      userID,
		BookshelfID: "testbookshelf",
		ISBN:        "1234567890",
		Title:       "Test Book",
		Author:      "Test Author",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	bookshelfRepo.On("GetByID", ctx, existingBook.BookshelfID).Return(
		&models.Bookshelf{
			ID:     existingBook.BookshelfID,
			UserID: userID,
		}, nil,
	)

	repo.On("GetByID", ctx, bookID).Return(existingBook, nil)
	repo.On("Delete", ctx, bookID).Return(nil)

	err := service.Delete(ctx, bookID)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
	bookshelfRepo.AssertExpectations(t)
}

func TestBookService_Delete_ErrorBookNotFound(t *testing.T) {
	repo := new(MockRepository)
	bookshelfRepo := new(MockBookshelfRepository)
	searcher := new(search.SearchService)
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	service := &BookService{
		repo:          repo,
		allBooksRepo:  repo,
		bookshelfRepo: bookshelfRepo,
		searcher:      searcher,
		log:           log,
		bookLimit:     1000,
	}

	ctx := context.Background()
	userID := "testuser"
	ctx = context.WithValue(ctx, "userID", userID)
	bookID := "nonexistentbookid"

	repo.On("GetByID", ctx, bookID).Return(nil, mongo.ErrBookNotFound)

	err := service.Delete(ctx, bookID)

	assert.ErrorIs(t, err, ErrBookNotFound)
	repo.AssertExpectations(t)
}
