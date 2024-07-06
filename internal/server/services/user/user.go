package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/getz-devs/librakeeper-server/internal/server/models"
	"github.com/getz-devs/librakeeper-server/internal/server/repository"
	"github.com/getz-devs/librakeeper-server/internal/server/services/book" // For ErrUserNotFoundInContext
	"github.com/getz-devs/librakeeper-server/internal/server/storage/mongo"
	"log/slog"
)

// Custom Error Type
var (
	ErrUserNotFoundInContext = errors.New("userID not found in context")
	ErrNotAuthorized         = errors.New("user is not authorized to perform this action")
	ErrUserNotFound          = errors.New("user not found")
)

// UserService handles business logic for user.
type UserService struct {
	repo repository.UserRepo
	log  *slog.Logger
}

// NewUserService creates a new UserService instance.
func NewUserService(repo repository.UserRepo, log *slog.Logger) *UserService {
	return &UserService{
		repo: repo,
		log:  log,
	}
}

// Create a new user.
func (s *UserService) Create(ctx context.Context, user *models.User) error {
	// Rule: Unique User ID
	exists, err := s.repo.ExistsByID(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
		return mongo.ErrUserAlreadyExists
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByID retrieves a user by ID.
func (s *UserService) GetByID(ctx context.Context, userID string) (*models.User, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, mongo.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// Update updates an existing user.
func (s *UserService) Update(ctx context.Context, userID string, update *models.UserUpdate) error {
	// Get userID from context (for authorization)
	ctxUserID, ok := ctx.Value("userID").(string)
	if !ok {
		return ErrUserNotFoundInContext
	}

	// Rule: User Self-Modification
	if userID != ctxUserID {
		return ErrNotAuthorized
	}

	if err := s.repo.Update(ctx, userID, update); err != nil {
		if errors.Is(err, mongo.ErrUserNotFound) {
			return ErrUserNotFound
		}
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// Delete deletes a user.
func (s *UserService) Delete(ctx context.Context, userID string) error {
	// Get userID from context (for authorization)
	ctxUserID, ok := ctx.Value("userID").(string)
	if !ok {
		return book.ErrUserNotFoundInContext
	}

	// Rule: User Self-Deletion
	if userID != ctxUserID {
		return ErrNotAuthorized
	}

	if err := s.repo.Delete(ctx, userID); err != nil {
		if errors.Is(err, mongo.ErrUserNotFound) {
			return ErrUserNotFound
		}
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
