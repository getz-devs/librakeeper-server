package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/getz-devs/librakeeper-server/internal/server/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
	"time"
)

var (
	ErrBookshelfNotFound = errors.New("bookshelf not found")
)

type BookshelfService struct {
	collection *mongo.Collection
	log        *slog.Logger
}

func NewBookshelfService(collection *mongo.Collection, log *slog.Logger) *BookshelfService {
	return &BookshelfService{
		collection: collection,
		log:        log,
	}
}

func (s *BookshelfService) CreateBookshelf(ctx context.Context, bookshelf *models.Bookshelf) (*models.Bookshelf, error) {
	bookshelf.CreatedAt = time.Now()
	bookshelf.UpdatedAt = time.Now()

	res, err := s.collection.InsertOne(ctx, bookshelf)
	if err != nil {
		s.log.Error("failed to create bookshelf", slog.Any("error", err))
		return nil, fmt.Errorf("failed to create bookshelf: %w", err)
	}

	bookshelf.ID = res.InsertedID.(primitive.ObjectID)
	return bookshelf, nil
}

func (s *BookshelfService) GetBookshelf(ctx context.Context, bookshelfID primitive.ObjectID) (*models.Bookshelf, error) {
	var bookshelf models.Bookshelf
	err := s.collection.FindOne(ctx, bson.M{"_id": bookshelfID}).Decode(&bookshelf)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrBookshelfNotFound
		}
		return nil, fmt.Errorf("failed to get bookshelf: %w", err)
	}
	return &bookshelf, nil
}

func (s *BookshelfService) GetBookshelvesByUserID(ctx context.Context, userID primitive.ObjectID, page int64, limit int64) ([]*models.Bookshelf, error) {
	findOptions := options.Find()
	findOptions.SetSkip((page - 1) * limit)
	findOptions.SetLimit(limit)

	cursor, err := s.collection.Find(ctx, bson.M{"user_id": userID}, findOptions)
	if err != nil {
		s.log.Error("failed to get bookshelves", slog.Any("error", err))
		return nil, fmt.Errorf("failed to get bookshelves: %w", err)
	}
	defer cursor.Close(ctx)

	var bookshelves []*models.Bookshelf
	for cursor.Next(ctx) {
		var bookshelf models.Bookshelf
		if err := cursor.Decode(&bookshelf); err != nil {
			return nil, fmt.Errorf("failed to decode bookshelf: %w", err)
		}
		bookshelves = append(bookshelves, &bookshelf)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}
	return bookshelves, nil
}

func (s *BookshelfService) UpdateBookshelf(ctx context.Context, bookshelfID primitive.ObjectID, update bson.M) error {
	update["updated_at"] = time.Now()
	_, err := s.collection.UpdateOne(ctx, bson.M{"_id": bookshelfID}, bson.M{"$set": update})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrBookshelfNotFound
		}
		return fmt.Errorf("failed to update bookshelf: %w", err)
	}
	return nil
}

func (s *BookshelfService) DeleteBookshelf(ctx context.Context, bookshelfID primitive.ObjectID) error {
	_, err := s.collection.DeleteOne(ctx, bson.M{"_id": bookshelfID})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrBookshelfNotFound
		}
		return fmt.Errorf("failed to delete bookshelf: %w", err)
	}
	return nil
}
