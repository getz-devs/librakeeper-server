package searcher_service

import (
	"context"
	"github.com/getz-devs/librakeeper-server/internal/searcher-shared/domain/bookModels"
	"log/slog"
)

type RequestExecutor interface {
	AddRequest(ctx context.Context, isbn string) error
}

type SearcherService struct {
	log             *slog.Logger
	requestStorage  RequestStorage
	requestExecutor RequestExecutor
}

type RequestStorage interface {
	FindOrCreateRequest(ctx context.Context, isbn string) (bookModels.SearchRequest, bool, error)
}

func New(log *slog.Logger, requestStorage RequestStorage, requestExecutor RequestExecutor) *SearcherService {
	return &SearcherService{
		log:             log,
		requestStorage:  requestStorage,
		requestExecutor: requestExecutor,
	}
}

func (s *SearcherService) SearchByISBN(ctx context.Context, isbn string) (bookModels.SearchRequest, error) {
	const op = "searcher.SearcherService.SearchByISBN"
	s.log.With(
		slog.String("op", op),
		slog.String("isbn", isbn),
	)
	s.log.Info("searching by ISBN")

	data, created, err := s.requestStorage.FindOrCreateRequest(ctx, isbn)
	if err != nil {
		return bookModels.SearchRequest{}, err
	}
	if created {
		s.log.Info("request created")
		err := s.requestExecutor.AddRequest(ctx, isbn)
		if err != nil {
			return bookModels.SearchRequest{}, err
		}
	}

	return data, nil
}

//func (s *SearcherService) FindOrCreateRequest(ctx context.Context, isbn string) (*models.BooksSearchResult, error) {
//	const op = "searcher.SearcherService.FindOrCreateRequest"
//	s.log.With(
//		slog.String("op", op),
//		slog.String("isbn", isbn),
//	)
//	s.log.Info("searching by ISBN")
//	return &models.BooksSearchResult{}, nil
//}
