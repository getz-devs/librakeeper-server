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
		collection: db.Collection("bookshelf"),
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
	var bookshelf models.Bookshelf
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&bookshelf)
	fmt.Printf("%+v\n", bookshelf)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrBookshelfNotFound
		}
		return nil, fmt.Errorf("failed to get bookshelf: %w", err)
	}

	return &bookshelf, nil
}

// GetByUser retrieves bookshelf associated with a specific user ID.
func (r *BookshelfRepo) GetByUser(ctx context.Context, userID string, page int64, limit int64) ([]*models.Bookshelf, error) {
	findOptions := options.Find()
	findOptions.SetSkip((page - 1) * limit)
	findOptions.SetLimit(limit)

	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID}, findOptions)
	if err != nil {
		r.log.Error("failed to get bookshelf by user ID", slog.Any("error", err))
		return nil, fmt.Errorf("failed to get bookshelf by user ID: %w", err)
	}
	defer cursor.Close(ctx)

	var bookshelves []*models.Bookshelf
	if err = cursor.All(ctx, &bookshelves); err != nil {
		return nil, fmt.Errorf("failed to decode bookshelf: %w", err)
	}

	return bookshelves, nil
}

// CountByUser returns the number of bookshelf owned by a user.
func (r *BookshelfRepo) CountByUser(ctx context.Context, userID string) (int, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"user_id": userID})
	if err != nil {
		r.log.Error("failed to count bookshelf by user ID", slog.Any("error", err))
		return 0, fmt.Errorf("failed to count bookshelf by user ID: %w", err)
	}

	return int(count), nil
}

// ExistsByNameAndUser checks if a bookshelf with the given name already exists for a user.
func (r *BookshelfRepo) ExistsByNameAndUser(ctx context.Context, name, userID string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"name": name, "user_id": userID})
	if err != nil {
		r.log.Error("failed to check bookshelf existence by name and user", slog.Any("error", err))
		return false, fmt.Errorf("failed to check bookshelf existence by name and user: %w", err)
	}

	return count > 0, nil
}

// Update updates a bookshelf in the database.
func (r *BookshelfRepo) Update(ctx context.Context, id string, update *models.BookshelfUpdate) error {
	update.UpdatedAt = time.Now()
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
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
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrBookshelfNotFound
		}
		return fmt.Errorf("failed to delete bookshelf: %w", err)
	}
	return nil
}
