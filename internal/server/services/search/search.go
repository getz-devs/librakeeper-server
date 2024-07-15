package search

import (
	"context"
	searcherv1 "github.com/getz-devs/librakeeper-protos/gen/go/searcher"
	"github.com/getz-devs/librakeeper-server/internal/server/repository"
	"log/slog"
)

// SearchService handles business logic for search.
type SearchService struct {
	searcher     repository.SearchRepo
	allBooksRepo repository.BookRepo
	log          *slog.Logger
}

func (s *SearchService) Simple(ctx context.Context, isbn string) (*searcherv1.SearchByISBNResponse, error) {
	// allBooksRepo.Search(isbn)
	return nil, nil
}

func (s *SearchService) Advanced(ctx context.Context, isbn string) (*searcherv1.SearchByISBNResponse, error) {
	// searcher.Search(isbn)
	return nil, nil
}

func NewSearchService(client repository.SearchRepo, repo repository.BookRepo, log *slog.Logger) *SearchService {
	return &SearchService{
		searcher:     client,
		allBooksRepo: repo,
	}
}
