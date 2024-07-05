package mongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/getz-devs/librakeeper-server/internal/server/models"
	"github.com/getz-devs/librakeeper-server/internal/server/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log/slog"
	"time"
)

// ErrUserNotFound occurs when a user is not found in the database.
var ErrUserNotFound = errors.New("user not found")

// ErrUserAlreadyExists occurs when trying to create a user with an ID that already exists.
var ErrUserAlreadyExists = errors.New("user already exists")

// UserRepo implements the repository.UserRepo interface for MongoDB.
type UserRepo struct {
	collection *mongo.Collection
	log        *slog.Logger
}

// NewUserRepo creates a new UserRepo instance.
func NewUserRepo(db *mongo.Database, log *slog.Logger) repository.UserRepo {
	return &UserRepo{
		collection: db.Collection("users"),
		log:        log,
	}
}

// Create inserts a new user into the database.
func (r *UserRepo) Create(ctx context.Context, user *models.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		// Check for duplicate key error
		var writeErr mongo.WriteException
		if errors.As(err, &writeErr) && writeErr.WriteErrors[0].Code == 11000 {
			return ErrUserAlreadyExists
		}

		r.log.Error("failed to create user", slog.Any("error", err))
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByID retrieves a user from the database by their ID.
func (r *UserRepo) GetByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// Update updates a user in the database.
func (r *UserRepo) Update(ctx context.Context, id string, update *models.UserUpdate) error {
	update.UpdatedAt = time.Now()
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrUserNotFound
		}
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// Delete removes a user from the database.
func (r *UserRepo) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrUserNotFound
		}
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
