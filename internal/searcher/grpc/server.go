package searcher

import (
	"context"
	searcherv1 "github.com/getz-devs/librakeeper-protos/gen/go/searcher"
	searcher_service "github.com/getz-devs/librakeeper-server/internal/searcher/services/searcher"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	searcherv1.UnimplementedSearcherServer
	searcherService *searcher_service.SearcherService
}

func Register(gRPC *grpc.Server, searcherService *searcher_service.SearcherService) {
	searcherv1.RegisterSearcherServer(gRPC, &serverAPI{searcherService: searcherService})
}

func (*serverAPI) SearchByISBN(
	ctx context.Context,
	req *searcherv1.SearchByISBNRequest,
) (*searcherv1.SearchByISBNResponse, error) {
	if req.GetIsbn() == "" {
		return nil, status.Error(codes.InvalidArgument, "isbn cannot be empty")
	}

	return &searcherv1.SearchByISBNResponse{
		Title:     "title",
		Author:    "author",
		Publisher: "publisher",
		Year:      "year",
	}, nil
}
