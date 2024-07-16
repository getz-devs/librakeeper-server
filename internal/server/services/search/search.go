package search

import (
	"context"
	"errors"
	searcherv1 "github.com/getz-devs/librakeeper-protos/gen/go/searcher"
	"github.com/getz-devs/librakeeper-server/internal/server/models"
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
func (s *SearchService) Simple(ctx context.Context, isbn string) (*models.SearchResponse, error) {
	const op = "search.SearchService.Simple"
	log := s.log.With(slog.String("op", op), slog.String("isbn", isbn))

	if isbn == "" {
		return nil, ErrISBNRequired
	}

	// Используем GetByISBN для поиска по ISBN
	book, err := s.allBooksRepo.GetByISBN(ctx, isbn)
	if err != nil {
		if errors.Is(err, mongo.ErrBookNotFound) {
			return nil, ErrISBNNotFound
		}
		log.Error("failed to get book from allBooksRepo", slog.Any("error", err))
		return nil, err
	}

	return &models.SearchResponse{
		Status: searcherv1.SearchByISBNResponse_SUCCESS,
		Books:  []*models.Book{book}, // Используем найденный объект book напрямую
	}, nil
}

// Advanced выполняет расширенный поиск по ISBN с использованием gRPC.
func (s *SearchService) Advanced(ctx context.Context, isbn string) (*models.SearchResponse, error) { // Измените тип возвращаемого значения
	const op = "search.SearchService.Advanced"
	log := s.log.With(slog.String("op", op), slog.String("isbn", isbn))

	if isbn == "" {
		return nil, ErrISBNRequired
	}

	grpcResponse, err := s.searcher.SearchByISBN(ctx, isbn) // Получаем ответ от gRPC сервиса
	if err != nil {
		log.Error("failed to search by ISBN", slog.Any("error", err))
		return nil, err
	}

	// Преобразуем gRPC ответ в models.SearchResponse
	var books []*models.Book
	for _, protoBook := range grpcResponse.Books {
		book := &models.Book{
			ISBN:       isbn,
			Title:      protoBook.Title,
			Author:     protoBook.Author,
			Publishing: protoBook.Publishing,
			CoverImage: protoBook.ImgUrl,
			ShopName:   protoBook.ShopName,
			// ... другие поля, если необходимо
		}
		books = append(books, book)
	}

	return &models.SearchResponse{
		Status: grpcResponse.Status,
		Books:  books,
	}, nil
}

func NewSearchService(client repository.SearchRepo, repo repository.BookRepo, log *slog.Logger) *SearchService {
	return &SearchService{
		searcher:     client,
		allBooksRepo: repo,
		log:          log,
	}
}
