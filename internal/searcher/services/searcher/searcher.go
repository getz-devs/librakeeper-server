package searcher_service

import (
	"context"
	"fmt"
	"github.com/getz-devs/librakeeper-server/internal/searcher/domain/models"
	"log/slog"
)

type SearcherService struct {
	log          *slog.Logger
	ISBNSearcher ISBNSearcher
}

// ISBNSearcher return models.BooksSearchResult
type ISBNSearcher interface {
	SearchByISBN(ctx context.Context, isbn string) (*models.BooksSearchResult, error)
}

func New(log *slog.Logger, ISBNSearcher ISBNSearcher) *SearcherService {
	return &SearcherService{
		log:          log,
		ISBNSearcher: ISBNSearcher,
	}
}

func (s *SearcherService) SearchByISBN(ctx context.Context, isbn string) (*models.BooksSearchResult, error) {
	const op = "searcher.SearcherService.SearchByISBN"
	s.log.With(
		slog.String("op", op),
		slog.String("isbn", isbn),
	)
	s.log.Info("searching by ISBN")

	data, err := s.ISBNSearcher.SearchByISBN(ctx, isbn)
	if err != nil {
		s.log.Error("error when searching by ISBN", err)
		return nil, fmt.Errorf("%s ,error when searching by ISBN: %w", op, err)
	}

	return data, nil
}
