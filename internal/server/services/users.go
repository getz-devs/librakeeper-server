package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/getz-devs/librakeeper-server/internal/server/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log/slog"
	"time"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)

type UserService struct {
	collection *mongo.Collection
	log        *slog.Logger
}

func NewUserService(collection *mongo.Collection, log *slog.Logger) *UserService {
	return &UserService{
		collection: collection,
		log:        log,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := s.collection.InsertOne(ctx, user)
	if err != nil {
		// Check for duplicate key error (Firebase UID uniqueness)
		var mongoErr mongo.WriteException
		if errors.As(err, &mongoErr) && mongoErr.WriteErrors[0].Code == 11000 {
			return nil, ErrUserAlreadyExists
		}

		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *UserService) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	var user models.User
	err := s.collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return &user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, userID string, update bson.M) error {
	update["updated_at"] = time.Now()
	_, err := s.collection.UpdateOne(ctx, bson.M{"_id": userID}, bson.M{"$set": update})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrUserNotFound
		}
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (s *UserService) DeleteUser(ctx context.Context, userID string) error {
	_, err := s.collection.DeleteOne(ctx, bson.M{"_id": userID})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrUserNotFound
		}
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
