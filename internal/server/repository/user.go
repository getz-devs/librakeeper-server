package repository

import (
	"context"
	"github.com/getz-devs/librakeeper-server/internal/server/models"
)

// UserRepo defines the interface for user repository operations.
type UserRepo interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id string) (*models.User, error)
	Update(ctx context.Context, id string, update *models.UserUpdate) error
	Delete(ctx context.Context, id string) error
}
