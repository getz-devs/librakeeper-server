package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/getz-devs/librakeeper-server/internal/server/config"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	_db *mongo.Database
)

func Initialize(cfg *config.Config, log *slog.Logger) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Database.URI))
	if err != nil {
		return nil, errors.New("failed to connect to mongodb")
	}

	_db = client.Database(cfg.Database.Name)
	log.Info("connected to mongodb", slog.String("database", cfg.Database.Name))

	return _db, nil
}

// Ping checks the database connectivity.
func Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second) // Use the context's timeout
	defer cancel()

	if _db == nil {
		return errors.New("mongodb has not been initialized")
	}

	if err := _db.Client().Ping(ctx, nil); err != nil {
		return fmt.Errorf("mongodb ping failed: %w", err)
	}
	return nil
}
