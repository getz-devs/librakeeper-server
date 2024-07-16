package search

import (
	"context"
	"errors"
	searcherv1 "github.com/getz-devs/librakeeper-protos/gen/go/searcher"
	"github.com/getz-devs/librakeeper-server/internal/server/repository"
	"github.com/getz-devs/librakeeper-server/internal/server/storage/mongo"
	"log/slog"
)

type SearchService struct {
	searcher     repository.SearchRepo
	allBooksRepo repository.BookRepo
	log          *slog.Logger
}

// Simple выполняет простой поиск по ISBN в локальной базе данных.
func (s *SearchService) Simple(ctx context.Context, isbn string) (*searcherv1.SearchByISBNResponse, error) {
	const op = "search.SearchService.Simple"
	log := s.log.With(slog.String("op", op), slog.String("isbn", isbn))

	if isbn == "" {
		return nil, ErrISBNRequired
	}

	book, err := s.allBooksRepo.GetByID(ctx, isbn)
	if err != nil {
		if errors.Is(err, mongo.ErrBookNotFound) {
			return nil, ErrISBNNotFound
		}
		log.Error("failed to get book from allBooksRepo", slog.Any("error", err))
		return nil, err
	}

	// Преобразуем models.Book в searcherv1.Book
	protoBook := &searcherv1.Book{
		Title:      book.Title,
		Author:     book.Author,
		Publishing: book.Description,
		ImgUrl:     book.CoverImage,
		ShopName:   "", // TODO: change book to proto book
	}

	return &searcherv1.SearchByISBNResponse{
		Status: searcherv1.SearchByISBNResponse_SUCCESS,
		Books:  []*searcherv1.Book{protoBook},
	}, nil
}

// Advanced выполняет расширенный поиск по ISBN с использованием gRPC.
func (s *SearchService) Advanced(ctx context.Context, isbn string) (*searcherv1.SearchByISBNResponse, error) {
	const op = "search.SearchService.Advanced"
	log := s.log.With(slog.String("op", op), slog.String("isbn", isbn))

	if isbn == "" {
		return nil, ErrISBNRequired
	}

	response, err := s.searcher.SearchByISBN(ctx, isbn)
	if err != nil {
		log.Error("failed to search by ISBN", slog.Any("error", err))
		return nil, err
	}

	// Здесь вы можете добавить дополнительную обработку ответа от gRPC-сервиса,
	// например, обогатить данные из allBooksRepo.

	return response, nil
}

func NewSearchService(client repository.SearchRepo, repo repository.BookRepo, log *slog.Logger) *SearchService {
	return &SearchService{
		searcher:     client,
		allBooksRepo: repo,
		log:          log,
	}
}
