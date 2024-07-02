package mongostorage

import (
	"context"
	"github.com/getz-devs/librakeeper-server/internal/searcher-shared/domain/bookModels"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

//func (s *Storage) SearchByISBN(ctx context.Context, isbn string) (*models.BooksSearchResult, error) {
//	const op = "storage.mongo.FindByISBN"
//
//	// Pass these options to the Find method
//	findOptions := options.Find()
//	findOptions.SetLimit(5)
//
//	// Here's an array in which you can store the decoded documents
//	var results models.BooksSearchResult
//
//	// Passing bson.D{{}} as the filter matches all documents in the collection
//	cur, err := s.col.Find(ctx, bson.D{{}}, findOptions)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Finding multiple documents returns a cursor
//	// Iterating through the cursor allows us to decode documents one at a time
//	for cur.Next(ctx) {
//
//		// create a value into which the single document can be decoded
//		var elem models.BookSearchResult
//		err := cur.Decode(&elem)
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		results = append(results, &elem)
//	}
//
//	if err := cur.Err(); err != nil {
//		return nil, fmt.Errorf("%s ,error when reading cursor: %w", op, err)
//	}
//
//	// Close the cursor once finished
//	if err := cur.Close(ctx); err != nil {
//		return nil, fmt.Errorf("%s ,error when closing cursor: %w", op, err)
//	}
//
//	return &results, nil
//}

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
