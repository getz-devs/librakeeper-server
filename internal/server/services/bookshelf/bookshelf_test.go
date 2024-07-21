package bookshelf

import (
	"context"
	"github.com/getz-devs/librakeeper-server/internal/server/models"
	"github.com/getz-devs/librakeeper-server/internal/server/storage/mongo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"os"
	"testing"
	"time"
)

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

func TestBookshelfService_Create_Success(t *testing.T) {
	repo := new(MockBookshelfRepository)
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	service := &BookshelfService{
		repo: repo,
		log:  log,
	}

	ctx := context.Background()
	userID := "testuser"
	ctx = context.WithValue(ctx, "userID", userID)

	bookshelf := &models.Bookshelf{
		Name: "Test Bookshelf",
	}

	repo.On("ExistsByNameAndUser", ctx, bookshelf.Name, userID).Return(false, nil)
	repo.On("Create", ctx, bookshelf).Return(nil)

	err := service.Create(ctx, bookshelf)
	assert.NoError(t, err)
	assert.Equal(t, userID, bookshelf.UserID)

	repo.AssertExpectations(t)
}

func TestBookshelfService_Create_ErrorNameRequired(t *testing.T) {
	repo := new(MockBookshelfRepository)
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	service := &BookshelfService{
		repo: repo,
		log:  log,
	}

	ctx := context.Background()
	userID := "testuser"
	ctx = context.WithValue(ctx, "userID", userID)

	bookshelf := &models.Bookshelf{}

	err := service.Create(ctx, bookshelf)

	assert.ErrorIs(t, err, ErrNameRequired)
	repo.AssertExpectations(t)
}

func TestBookshelfService_Create_ErrorUserNotFoundInContext(t *testing.T) {
	repo := new(MockBookshelfRepository)
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	service := &BookshelfService{
		repo: repo,
		log:  log,
	}

	ctx := context.Background()

	bookshelf := &models.Bookshelf{
		Name: "Test Bookshelf",
	}

	err := service.Create(ctx, bookshelf)

	assert.ErrorIs(t, err, ErrUserNotFoundInContext)
	repo.AssertExpectations(t)
}

func TestBookshelfService_GetByID_Success(t *testing.T) {
	repo := new(MockBookshelfRepository)
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	service := &BookshelfService{
		repo: repo,
		log:  log,
	}

	ctx := context.Background()
	bookshelfID := "testbookshelfid"

	expectedBookshelf := &models.Bookshelf{
		ID:        bookshelfID,
		UserID:    "testuser",
		Name:      "Test Bookshelf",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repo.On("GetByID", ctx, bookshelfID).Return(expectedBookshelf, nil)

	bookshelf, err := service.GetByID(ctx, bookshelfID)

	assert.NoError(t, err)
	assert.Equal(t, expectedBookshelf, bookshelf)
	repo.AssertExpectations(t)
}

func TestBookshelfService_GetByID_ErrorBookshelfNotFound(t *testing.T) {
	repo := new(MockBookshelfRepository)
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	service := &BookshelfService{
		repo: repo,
		log:  log,
	}

	ctx := context.Background()
	bookshelfID := "nonexistentbookshelf"

	repo.On("GetByID", ctx, bookshelfID).Return(nil, mongo.ErrBookshelfNotFound)

	bookshelf, err := service.GetByID(ctx, bookshelfID)

	assert.ErrorIs(t, err, ErrBookshelfNotFound)
	assert.Nil(t, bookshelf)
	repo.AssertExpectations(t)
}

func TestBookshelfService_Update_Success(t *testing.T) {
	repo := new(MockBookshelfRepository)
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	service := &BookshelfService{
		repo: repo,
		log:  log,
	}

	ctx := context.Background()
	userID := "testuser"
	ctx = context.WithValue(ctx, "userID", userID)
	bookshelfID := "testbookshelfid"
	update := &models.BookshelfUpdate{
		Name: stringPtr("Updated Bookshelf Name"),
	}

	existingBookshelf := &models.Bookshelf{
		ID:        bookshelfID,
		UserID:    userID,
		Name:      "Test Bookshelf",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repo.On("GetByID", ctx, bookshelfID).Return(existingBookshelf, nil)
	repo.On("Update", ctx, bookshelfID, update).Return(nil)

	err := service.Update(ctx, bookshelfID, update)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestBookshelfService_Update_ErrorBookshelfNotFound(t *testing.T) {
	repo := new(MockBookshelfRepository)
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	service := &BookshelfService{
		repo: repo,
		log:  log,
	}

	ctx := context.Background()
	userID := "testuser"
	ctx = context.WithValue(ctx, "userID", userID)
	bookshelfID := "nonexistentbookshelf"
	update := &models.BookshelfUpdate{
		Name: stringPtr("Updated Bookshelf Name"),
	}

	repo.On("GetByID", ctx, bookshelfID).Return(nil, mongo.ErrBookshelfNotFound)

	err := service.Update(ctx, bookshelfID, update)

	assert.ErrorIs(t, err, ErrBookshelfNotFound)
	repo.AssertExpectations(t)
}

func TestBookshelfService_Delete_Success(t *testing.T) {
	repo := new(MockBookshelfRepository)
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	service := &BookshelfService{
		repo: repo,
		log:  log,
	}

	ctx := context.Background()
	userID := "testuser"
	ctx = context.WithValue(ctx, "userID", userID)
	bookshelfID := "testbookshelfid"

	existingBookshelf := &models.Bookshelf{
		ID:        bookshelfID,
		UserID:    userID,
		Name:      "Test Bookshelf",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repo.On("GetByID", ctx, bookshelfID).Return(existingBookshelf, nil)
	repo.On("Delete", ctx, bookshelfID).Return(nil)

	err := service.Delete(ctx, bookshelfID)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestBookshelfService_Delete_ErrorBookshelfNotFound(t *testing.T) {
	repo := new(MockBookshelfRepository)
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	service := &BookshelfService{
		repo: repo,
		log:  log,
	}

	ctx := context.Background()
	userID := "testuser"
	ctx = context.WithValue(ctx, "userID", userID)
	bookshelfID := "nonexistentbookshelf"

	repo.On("GetByID", ctx, bookshelfID).Return(nil, mongo.ErrBookshelfNotFound)

	err := service.Delete(ctx, bookshelfID)

	assert.ErrorIs(t, err, ErrBookshelfNotFound)
	repo.AssertExpectations(t)
}

// Helper function to create a pointer to a string
func stringPtr(s string) *string {
	return &s
}
