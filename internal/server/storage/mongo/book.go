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

// ErrBookNotFound occurs when a book is not found in the database.
var ErrBookNotFound = errors.New("book not found")

// ErrBookAlreadyExists occurs when trying to create a book with an ID that already exists.
var ErrBookAlreadyExists = errors.New("book already exists")

// BookRepo implements the repository.BookRepo interface for MongoDB.
type BookRepo struct {
	collection *mongo.Collection
	log        *slog.Logger
}

// NewBookRepo creates a new BookRepo instance.
func NewBookRepo(db *mongo.Database, log *slog.Logger, collectionName string) repository.BookRepo {
	return &BookRepo{
		collection: db.Collection(collectionName),
		log:        log,
	}
}

// Create inserts a new book into the database.
func (r *BookRepo) Create(ctx context.Context, book *models.Book) error {
	book.ID = primitive.NewObjectID().Hex()
	book.CreatedAt = time.Now()
	book.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, book)
	if err != nil {
		// Check for duplicate key error
		var writeErr mongo.WriteException
		if errors.As(err, &writeErr) && writeErr.WriteErrors[0].Code == 11000 {
			return ErrBookAlreadyExists
		}

		r.log.Error("failed to create book", slog.Any("error", err))
		return fmt.Errorf("failed to create book: %w", err)
	}

	return nil
}

// GetByID retrieves a book from the database by its ID.
func (r *BookRepo) GetByID(ctx context.Context, id string) (*models.Book, error) {
	var book models.Book
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&book)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrBookNotFound
		}
		return nil, fmt.Errorf("failed to get book: %w", err)
	}
	return &book, nil
}

// GetByISBN retrieves a book from the database by its ISBN.
func (r *BookRepo) GetByISBN(ctx context.Context, isbn string) (*models.Book, error) {
	var book models.Book
	err := r.collection.FindOne(ctx, bson.M{"isbn": isbn}).Decode(&book)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrBookNotFound
		}
		return nil, fmt.Errorf("failed to get book by ISBN: %w", err)
	}
	return &book, nil
}

// GetByUserID retrieves book associated with a specific user ID.
func (r *BookRepo) GetByUserID(ctx context.Context, userID string, page int64, limit int64) ([]*models.Book, error) {
	findOptions := options.Find()
	findOptions.SetSkip((page - 1) * limit)
	findOptions.SetLimit(limit)

	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID}, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to get book by user ID: %w", err)
	}
	defer cursor.Close(ctx)

	var books []*models.Book
	if err = cursor.All(ctx, &books); err != nil {
		return nil, fmt.Errorf("failed to decode book: %w", err)
	}

	return books, nil
}

// GetByBookshelfID retrieves book belonging to a specific bookshelf ID.
func (r *BookRepo) GetByBookshelfID(ctx context.Context, bookshelfID string, page int64, limit int64) ([]*models.Book, error) {
	matchStage := bson.D{{"$match", bson.D{{"bookshelf_id", bookshelfID}}}}
	skipStage := bson.D{{"$skip", (page - 1) * limit}}
	limitStage := bson.D{{"$limit", limit}}

	cursor, err := r.collection.Aggregate(ctx, mongo.Pipeline{matchStage, skipStage, limitStage})
	if err != nil {
		r.log.Error("failed to get book by bookshelf id", slog.Any("error", err))
		return nil, fmt.Errorf("failed to get book by bookshelf id: %w", err)
	}
	defer cursor.Close(ctx)

	var books []*models.Book
	if err = cursor.All(ctx, &books); err != nil {
		return nil, fmt.Errorf("failed to decode book: %w", err)
	}

	return books, nil
}

// CountInBookshelf returns the number of book in a bookshelf.
func (r *BookRepo) CountInBookshelf(ctx context.Context, bookshelfID string) (int, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"bookshelf_id": bookshelfID})
	if err != nil {
		r.log.Error("failed to count book by bookshelf ID", slog.Any("error", err))
		return 0, fmt.Errorf("failed to count book by bookshelf ID: %w", err)
	}

	return int(count), nil
}

// ExistsInBookshelf checks if a book with the given ISBN already exists in the bookshelf.
func (r *BookRepo) ExistsInBookshelf(ctx context.Context, isbn, bookshelfID string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"isbn": isbn, "bookshelf_id": bookshelfID})
	if err != nil {
		r.log.Error("failed to check book existence", slog.Any("error", err))
		return false, fmt.Errorf("failed to check book existence: %w", err)
	}

	return count > 0, nil
}

// Update updates a book in the database.
func (r *BookRepo) Update(ctx context.Context, id string, update *models.BookUpdate) error {
	update.UpdatedAt = time.Now()
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrBookNotFound
		}
		return fmt.Errorf("failed to update book: %w", err)
	}
	return nil
}

// Delete removes a book from the database.
func (r *BookRepo) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrBookNotFound
		}
		return fmt.Errorf("failed to delete book: %w", err)
	}
	return nil
}
