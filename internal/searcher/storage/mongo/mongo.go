package mongostorage

import (
	"context"
	"fmt"
	"github.com/getz-devs/librakeeper-server/internal/searcher/domain/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
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
	const op = "storage.mongo.New"

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
	const op = "storage.mongo.Close"
	if err := s.client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}

func (s *Storage) SearchByISBN(ctx context.Context, isbn string) (*models.BooksSearchResult, error) {
	const op = "storage.mongo.FindByISBN"

	// Pass these options to the Find method
	findOptions := options.Find()
	findOptions.SetLimit(5)

	// Here's an array in which you can store the decoded documents
	var results models.BooksSearchResult

	// Passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := s.col.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem models.BookSearchResult
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("%s ,error when reading cursor: %w", op, err)
	}

	// Close the cursor once finished
	if err := cur.Close(ctx); err != nil {
		return nil, fmt.Errorf("%s ,error when closing cursor: %w", op, err)
	}

	return &results, nil
}
