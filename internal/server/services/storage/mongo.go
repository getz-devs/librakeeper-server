package storage

import (
	"context"
	"fmt"
	"github.com/getz-devs/librakeeper-server/internal/server/config"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Collections struct {
	UsersCollection       *mongo.Collection
	BooksCollection       *mongo.Collection
	BookshelvesCollection *mongo.Collection
}

var (
	_log         *slog.Logger
	_db          *mongo.Database
	_collections Collections
)

func Initialize(cfg *config.Config, log *slog.Logger) (*mongo.Database, Collections) {
	_log = log
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		_log.Error("failed to connect to mongodb", slog.Any("error", err))
		panic(err)
	}

	_db = client.Database(cfg.Database)
	_log.Info("connected to mongodb", slog.String("database", cfg.Database))

	_collections.UsersCollection = _db.Collection("users")
	_collections.BooksCollection = _db.Collection("books")
	_collections.BookshelvesCollection = _db.Collection("bookshelves")

	return _db, _collections
}

func GetCollection(name string) *mongo.Collection {
	if _db == nil {
		panic(fmt.Errorf("mongodb has not been initialized"))
	}

	return _db.Collection(name)
}

// Ping checks the database connectivity.
func Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second) // Use the context's timeout
	defer cancel()

	if _db == nil {
		panic(fmt.Errorf("mongodb has not been initialized"))
	}

	if err := _db.Client().Ping(ctx, nil); err != nil {
		return fmt.Errorf("mongodb ping failed: %w", err)
	}
	return nil
}
