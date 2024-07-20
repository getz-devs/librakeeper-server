package search

import (
	"context"
	"errors"
	searcherv1 "github.com/getz-devs/librakeeper-protos/gen/go/searcher"
	"github.com/getz-devs/librakeeper-server/internal/server/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"os"
	"testing"
	"time"
)

type MockSearchRepo struct {
	mock.Mock
}

func (m *MockSearchRepo) SearchByISBN(ctx context.Context, isbn string) (*searcherv1.SearchByISBNResponse, error) {
	args := m.Called(ctx, isbn)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*searcherv1.SearchByISBNResponse), args.Error(1)
}

type MockBookRepo struct {
	mock.Mock
}

func (m *MockBookRepo) Create(ctx context.Context, book *models.Book) error {
	args := m.Called(ctx, book)
	return args.Error(0)
}

func (m *MockBookRepo) GetByID(ctx context.Context, id string) (*models.Book, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Book), args.Error(1)
}

func (m *MockBookRepo) GetByISBN(ctx context.Context, isbn string) (*models.Book, error) {
	args := m.Called(ctx, isbn)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Book), args.Error(1)
}

func (m *MockBookRepo) GetByUserID(ctx context.Context, userID string, page int64, limit int64) ([]*models.Book, error) {
	args := m.Called(ctx, userID, page, limit)
	return args.Get(0).([]*models.Book), args.Error(1)
}

func (m *MockBookRepo) GetByBookshelfID(ctx context.Context, bookshelfID string, page int64, limit int64) ([]*models.Book, error) {
	args := m.Called(ctx, bookshelfID, page, limit)
	return args.Get(0).([]*models.Book), args.Error(1)
}

func (m *MockBookRepo) CountInBookshelf(ctx context.Context, bookshelfID string) (int, error) {
	args := m.Called(ctx, bookshelfID)
	return args.Int(0), args.Error(1)
}

func (m *MockBookRepo) ExistsInBookshelf(ctx context.Context, isbn, bookshelfID string) (bool, error) {
	args := m.Called(ctx, isbn, bookshelfID)
	return args.Bool(0), args.Error(1)
}

func (m *MockBookRepo) Update(ctx context.Context, id string, update *models.BookUpdate) error {
	args := m.Called(ctx, id, update)
	return args.Error(0)
}

func (m *MockBookRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestSearchService_Simple_Success(t *testing.T) {
	searchRepo := new(MockSearchRepo)
	bookRepo := new(MockBookRepo)
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	service := &SearchService{
		searcher:     searchRepo,
		allBooksRepo: bookRepo,
		log:          log,
	}

	ctx := context.Background()
	isbn := "1234567890"
	expectedBook := &models.Book{
		ID:          "testbookid",
		UserID:      "testuser",
		BookshelfID: "testbookshelf",
		ISBN:        isbn,
		Title:       "Test Book",
		Author:      "Test Author",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	bookRepo.On("GetByISBN", ctx, isbn).Return(expectedBook, nil)

	resp, err := service.Simple(ctx, isbn)

	assert.NoError(t, err)
	assert.Equal(t, searcherv1.SearchByISBNResponse_SUCCESS, resp.Status)
	assert.Equal(t, expectedBook, resp.Books[0])
	bookRepo.AssertExpectations(t)
}

func TestSearchService_Simple_ErrorISBNNotFound(t *testing.T) {
	searchRepo := new(MockSearchRepo)
	bookRepo := new(MockBookRepo)
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	service := &SearchService{
		searcher:     searchRepo,
		allBooksRepo: bookRepo,
		log:          log,
	}

	ctx := context.Background()
	isbn := "nonexistentisbn"

	bookRepo.On("GetByISBN", ctx, isbn).Return(nil, errors.New("isbn not found"))

	_, err := service.Simple(ctx, isbn)
	assert.ErrorContains(t, err, "isbn not found")
	bookRepo.AssertExpectations(t)
}

func TestSearchService_Simple_ErrorISBNRequired(t *testing.T) {
	searchRepo := new(MockSearchRepo)
	bookRepo := new(MockBookRepo)
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	service := &SearchService{
		searcher:     searchRepo,
		allBooksRepo: bookRepo,
		log:          log,
	}

	ctx := context.Background()
	isbn := ""

	_, err := service.Simple(ctx, isbn)
	assert.ErrorIs(t, err, ErrISBNRequired)
}

func TestSearchService_Advanced_Success(t *testing.T) {
	searchRepo := new(MockSearchRepo)
	bookRepo := new(MockBookRepo)
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	service := &SearchService{
		searcher:     searchRepo,
		allBooksRepo: bookRepo,
		log:          log,
	}

	ctx := context.Background()
	isbn := "1234567890"

	grpcResponse := &searcherv1.SearchByISBNResponse{
		Status: searcherv1.SearchByISBNResponse_SUCCESS,
		Books: []*searcherv1.Book{
			{
				Title:      "Test Book 1",
				Author:     "Test Author 1",
				Publishing: "Test Publishing 1",
				ImgUrl:     "https://example.com/cover1.jpg",
				ShopName:   "Test Shop 1",
			},
			{
				Title:      "Test Book 2",
				Author:     "Test Author 2",
				Publishing: "Test Publishing 2",
				ImgUrl:     "https://example.com/cover2.jpg",
				ShopName:   "Test Shop 2",
			},
		},
	}

	searchRepo.On("SearchByISBN", ctx, isbn).Return(grpcResponse, nil)

	resp, err := service.Advanced(ctx, isbn)

	assert.NoError(t, err)
	assert.Equal(t, grpcResponse.Status, resp.Status)
	assert.Len(t, resp.Books, len(grpcResponse.Books))
	searchRepo.AssertExpectations(t)
}

func TestSearchService_Advanced_ErrorISBNNotFound(t *testing.T) {
	searchRepo := new(MockSearchRepo)
	bookRepo := new(MockBookRepo)
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	service := &SearchService{
		searcher:     searchRepo,
		allBooksRepo: bookRepo,
		log:          log,
	}

	ctx := context.Background()
	isbn := "nonexistentisbn"

	searchRepo.On("SearchByISBN", ctx, isbn).Return(nil, errors.New("isbn not found"))

	_, err := service.Advanced(ctx, isbn)
	assert.ErrorContains(t, err, "isbn not found")
	searchRepo.AssertExpectations(t)
}

func TestSearchService_Advanced_ErrorISBNRequired(t *testing.T) {
	searchRepo := new(MockSearchRepo)
	bookRepo := new(MockBookRepo)
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	service := &SearchService{
		searcher:     searchRepo,
		allBooksRepo: bookRepo,
		log:          log,
	}

	ctx := context.Background()
	isbn := ""

	_, err := service.Advanced(ctx, isbn)
	assert.ErrorIs(t, err, ErrISBNRequired)
}
