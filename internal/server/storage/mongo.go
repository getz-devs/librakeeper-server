package database

import (
	"context"
	"fmt"
	"github.com/getz-devs/librakeeper-server/internal/server/config"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	_db  *mongo.Database
	_log *slog.Logger
)

func Initialize(cfg *config.Config, log *slog.Logger) *mongo.Database {
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
	return _db
}

func GetCollection(name string) *mongo.Collection {
	if _db == nil {
		panic(fmt.Errorf("mongodb has not been initialized"))
	}

	return _db.Collection(name)
}
