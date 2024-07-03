package mongostorage

import (
	"context"
	"github.com/getz-devs/librakeeper-server/internal/searcher-shared/domain/bookModels"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Storage struct {
	client *mongo.Client
	col    *mongo.Collection
}

type DatabaseMongoConfig struct {
	ConnectUrl string
	Database   string
	Collection string
}

func New(databaseConfig DatabaseMongoConfig) *Storage {
	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(databaseConfig.ConnectUrl))
	if err != nil {
		panic(err)
	}

	coll := client.Database(databaseConfig.Database).Collection(databaseConfig.Collection)

	return &Storage{
		client: client,
		col:    coll,
	}
}

func (s *Storage) Close() {
	if err := s.client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}

func (s *Storage) FindOrCreateRequest(ctx context.Context, isbn string) (bookModels.SearchRequest, bool, error) {
	// insert if not exist (upsert)
	filter := bson.D{{"isbn", isbn}}
	insertValue := bookModels.New(isbn)
	value := bson.D{{"$setOnInsert", insertValue}}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var result bookModels.SearchRequest
	err := s.col.FindOneAndUpdate(ctx, filter, value, opts).Decode(&result)
	if err != nil {
		return bookModels.SearchRequest{}, false, err
	}
	if insertValue.ID == result.ID {
		return bookModels.SearchRequest{}, true, nil
	}
	return result, false, nil
}

func (s *Storage) CompleteRequest(ctx context.Context, isbn string, books []*bookModels.BookInShop) error {
	filter := bson.D{{"isbn", isbn}}
	values := bson.D{{"$set", bson.D{
		{"books", books},
		{"status", bookModels.Success},
		{"updated_at", time.Now()},
	}}}
	if _, err := s.col.UpdateOne(ctx, filter, values); err != nil {
		return err
	}
	return nil
}

func (s *Storage) RejectRequest(ctx context.Context, isbn string) error {
	filter := bson.D{{"isbn", isbn}}
	values := bson.D{{"$set", bson.D{
		{"status", bookModels.Failed},
		{"updated_at", time.Now()},
	}}}
	if _, err := s.col.UpdateOne(ctx, filter, values); err != nil {
		return err
	}
	return nil
}
