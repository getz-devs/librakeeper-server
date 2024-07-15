package search

import (
	"context"
	"errors"
	searcherv1 "github.com/getz-devs/librakeeper-protos/gen/go/searcher"
	"github.com/getz-devs/librakeeper-server/internal/server/repository"
	"log/slog"
)

var (
	ErrISBNNotFound = errors.New("ISBN not found")
	ErrISBNRequired = errors.New("ISBN is required")
)

type SearcherClient struct {
	// grpc
	log *slog.Logger
}

func (s SearcherClient) SearchByISBN(ctx context.Context, isbn string) (*searcherv1.SearchByISBNResponse, error) {
	//TODO implement me
	panic("implement me")
	return nil, nil
}

func NewSearcherClient(log *slog.Logger) repository.SearchRepo {
	return &SearcherClient{
		log: log,
	}
}
