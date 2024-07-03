package searcher

import (
	"context"
	searcherv1 "github.com/getz-devs/librakeeper-protos/gen/go/searcher"
	"github.com/getz-devs/librakeeper-server/internal/searcher-shared/domain/bookModels"
	searcherservice "github.com/getz-devs/librakeeper-server/internal/searcher/services/searcher"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

type serverAPI struct {
	searcherv1.UnimplementedSearcherServer
	searcherService *searcherservice.SearcherService
	log             *slog.Logger
}

func Register(gRPC *grpc.Server, searcherService *searcherservice.SearcherService, log *slog.Logger) {
	searcherv1.RegisterSearcherServer(gRPC, &serverAPI{searcherService: searcherService, log: log})
}

func (s *serverAPI) SearchByISBN(
	ctx context.Context,
	req *searcherv1.SearchByISBNRequest,
) (*searcherv1.SearchByISBNResponse, error) {
	if req.GetIsbn() == "" {
		return nil, status.Error(codes.InvalidArgument, "isbn cannot be empty")
	}
	results, err := s.searcherService.SearchByISBN(ctx, req.GetIsbn())
	if err != nil {
		return nil, err
	}
	s.log.Info("Results", slog.Any("results", results))

	var books []*searcherv1.Book
	if len(results.Books) >= 0 {
		for _, book := range results.Books {
			books = append(books, &searcherv1.Book{
				Title:      book.Title,
				Author:     book.Author,
				Publishing: book.Publishing,
				ImgUrl:     book.ImgUrl,
				ShopName:   book.ShopName,
			})
		}
	}

	statusReturn := searcherv1.SearchByISBNResponse_PROCESSING
	switch results.Status {
	case bookModels.Pending:
		statusReturn = searcherv1.SearchByISBNResponse_PROCESSING
	case bookModels.Failed:
		statusReturn = searcherv1.SearchByISBNResponse_PROCESSING
	case bookModels.Success:
		statusReturn = searcherv1.SearchByISBNResponse_SUCCESS
	default:
		statusReturn = searcherv1.SearchByISBNResponse_PROCESSING
	}
	return &searcherv1.SearchByISBNResponse{
		Status: statusReturn,
		Books:  books,
	}, nil
}
