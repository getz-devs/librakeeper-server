package mongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/getz-devs/librakeeper-server/internal/server/models"
	"github.com/getz-devs/librakeeper-server/internal/server/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
	"time"
)

// ErrBookshelfNotFound occurs when a bookshelf is not found in the database.
var ErrBookshelfNotFound = errors.New("bookshelf not found")

// ErrBookshelfAlreadyExists occurs when trying to create a bookshelf with an ID that already exists.
var ErrBookshelfAlreadyExists = errors.New("bookshelf already exists")

// BookshelfRepo implements the repository.BookshelfRepo interface for MongoDB.
type BookshelfRepo struct {
	collection *mongo.Collection
	log        *slog.Logger
}

// NewBookshelfRepo creates a new BookshelfRepo instance.
func NewBookshelfRepo(db *mongo.Database, log *slog.Logger) repository.BookshelfRepo {
	return &BookshelfRepo{
		collection: db.Collection("bookshelves"),
		log:        log,
	}
}

// Create inserts a new bookshelf into the database.
func (r *BookshelfRepo) Create(ctx context.Context, bookshelf *models.Bookshelf) error {
	bookshelf.ID = primitive.NewObjectID().Hex()
	bookshelf.CreatedAt = time.Now()
	bookshelf.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, bookshelf)
	if err != nil {
		// Check for duplicate key error
		var writeErr mongo.WriteException
		if errors.As(err, &writeErr) && writeErr.WriteErrors[0].Code == 11000 {
			return ErrBookshelfAlreadyExists
		}

		r.log.Error("failed to create bookshelf", slog.Any("error", err))
		return fmt.Errorf("failed to create bookshelf: %w", err)
	}

	return nil
}

// GetByID retrieves a bookshelf from the database by its ID.
func (r *BookshelfRepo) GetByID(ctx context.Context, id string) (*models.Bookshelf, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid bookshelf ID: %w", err)
	}

	var bookshelf models.Bookshelf
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&bookshelf)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrBookshelfNotFound
		}
		return nil, fmt.Errorf("failed to get bookshelf: %w", err)
	}

	return &bookshelf, nil
}

// GetByUserID retrieves bookshelves associated with a specific user ID.
func (r *BookshelfRepo) GetByUserID(ctx context.Context, userID string, page int64, limit int64) ([]*models.Bookshelf, error) {
	objectUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	findOptions := options.Find()
	findOptions.SetSkip((page - 1) * limit)
	findOptions.SetLimit(limit)

	cursor, err := r.collection.Find(ctx, bson.M{"user_id": objectUserID}, findOptions)
	if err != nil {
		r.log.Error("failed to get bookshelves by user ID", slog.Any("error", err))
		return nil, fmt.Errorf("failed to get bookshelves by user ID: %w", err)
	}
	defer cursor.Close(ctx)

	var bookshelves []*models.Bookshelf
	if err = cursor.All(ctx, &bookshelves); err != nil {
		return nil, fmt.Errorf("failed to decode bookshelves: %w", err)
	}

	return bookshelves, nil
}

// Update updates a bookshelf in the database.
func (r *BookshelfRepo) Update(ctx context.Context, id string, update *models.BookshelfUpdate) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid bookshelf ID: %w", err)
	}

	update.UpdatedAt = time.Now()
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": update})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrBookshelfNotFound
		}
		return fmt.Errorf("failed to update bookshelf: %w", err)
	}
	return nil
}

// Delete removes a bookshelf from the database.
func (r *BookshelfRepo) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid bookshelf ID: %w", err)
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrBookshelfNotFound
		}
		return fmt.Errorf("failed to delete bookshelf: %w", err)
	}
	return nil
}
