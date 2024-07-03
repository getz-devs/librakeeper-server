package books

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
	ErrBookNotFound = errors.New("book not found")
)

type BookService struct {
	collection *mongo.Collection
	log        *slog.Logger
}

func NewBookService(collection *mongo.Collection, log *slog.Logger) *BookService {
	return &BookService{
		collection: collection,
		log:        log,
	}
}

func (s *BookService) CreateBook(ctx context.Context, book *models.Book) (*models.Book, error) {
	book.CreatedAt = time.Now()
	book.UpdatedAt = time.Now()

	res, err := s.collection.InsertOne(ctx, book)
	if err != nil {
		s.log.Error("failed to create book", slog.Any("error", err))
		return nil, fmt.Errorf("failed to create book: %w", err)
	}

	book.ID = res.InsertedID.(primitive.ObjectID)
	return book, nil
}

func (s *BookService) GetBook(ctx context.Context, bookID primitive.ObjectID) (*models.Book, error) {
	var book models.Book
	err := s.collection.FindOne(ctx, bson.M{"_id": bookID}).Decode(&book)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrBookNotFound
		}
		return nil, fmt.Errorf("failed to get book: %w", err)
	}
	return &book, nil
}

func (s *BookService) GetBooks(ctx context.Context, page int64, limit int64) ([]*models.Book, error) {
	findOptions := options.Find()
	findOptions.SetSkip((page - 1) * limit)
	findOptions.SetLimit(limit)

	cursor, err := s.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		s.log.Error("failed to get books", slog.Any("error", err))
		return nil, fmt.Errorf("failed to get books: %w", err)
	}
	defer cursor.Close(ctx)

	var books []*models.Book
	for cursor.Next(ctx) {
		var book models.Book
		if err := cursor.Decode(&book); err != nil {
			return nil, fmt.Errorf("failed to decode book: %w", err)
		}
		books = append(books, &book)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}
	return books, nil
}

func (s *BookService) GetBooksByBookshelfID(ctx context.Context, bookshelfID primitive.ObjectID, page int64, limit int64) ([]*models.Book, error) {
	matchStage := bson.D{{"$match", bson.D{{"bookshelf_id", bookshelfID}}}}
	skipStage := bson.D{{"$skip", (page - 1) * limit}}
	limitStage := bson.D{{"$limit", limit}}

	cursor, err := s.collection.Aggregate(ctx, mongo.Pipeline{matchStage, skipStage, limitStage})
	if err != nil {
		s.log.Error("failed to get books by bookshelf id", slog.Any("error", err))
		return nil, fmt.Errorf("failed to get books by bookshelf id: %w", err)
	}
	defer cursor.Close(ctx)

	var books []*models.Book
	for cursor.Next(ctx) {
		var book models.Book
		if err := cursor.Decode(&book); err != nil {
			return nil, fmt.Errorf("failed to decode book: %w", err)
		}
		books = append(books, &book)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}
	return books, nil
}

func (s *BookService) UpdateBook(ctx context.Context, bookID primitive.ObjectID, update bson.M) error {
	update["updated_at"] = time.Now()
	_, err := s.collection.UpdateOne(ctx, bson.M{"_id": bookID}, bson.M{"$set": update})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrBookNotFound
		}
		return fmt.Errorf("failed to update book: %w", err)
	}
	return nil
}

func (s *BookService) DeleteBook(ctx context.Context, bookID primitive.ObjectID) error {
	_, err := s.collection.DeleteOne(ctx, bson.M{"_id": bookID})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrBookNotFound
		}
		return fmt.Errorf("failed to delete book: %w", err)
	}
	return nil
}
