package search

import (
	"context"
	"errors"
	searcherv1 "github.com/getz-devs/librakeeper-protos/gen/go/searcher"
	"github.com/getz-devs/librakeeper-server/internal/server/repository"
	"google.golang.org/grpc"
	"log/slog"
	"time"
)

var (
	ErrISBNNotFound = errors.New("ISBN not found")
	ErrISBNRequired = errors.New("ISBN is required")
)

type SearcherClient struct {
	client searcherv1.SearcherClient // Используем gRPC клиент напрямую
	log    *slog.Logger
}

func (s *SearcherClient) SearchByISBN(ctx context.Context, isbn string) (*searcherv1.SearchByISBNResponse, error) {
	const op = "search.SearcherClient.SearchByISBN"
	log := s.log.With(slog.String("op", op), slog.String("isbn", isbn))

	if isbn == "" {
		return nil, ErrISBNRequired
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	// Используем client напрямую, без создания нового
	resp, err := s.client.SearchByISBN(ctx, &searcherv1.SearchByISBNRequest{Isbn: isbn})
	if err != nil {
		log.Error("failed to search by isbn", slog.Any("error", err))
		return nil, err
	}

	return resp, nil
}

func NewSearcherClient(conn *grpc.ClientConn, log *slog.Logger) repository.SearchRepo { // Передаем конфигурацию
	client := searcherv1.NewSearcherClient(conn) // Создаем gRPC клиент

	return &SearcherClient{
		client: client, // Сохраняем клиент в структуре
		log:    log,
	}
}
