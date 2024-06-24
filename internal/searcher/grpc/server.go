package searcher

import (
	"context"
	searcherv1 "github.com/getz-devs/librakeeper-protos/gen/go/searcher"
	"google.golang.org/grpc"
)

type serverAPI struct {
	searcherv1.UnimplementedSearcherServer
}

func Register(gRPC *grpc.Server) {
	searcherv1.RegisterSearcherServer(gRPC, &serverAPI{})
}

func (*serverAPI) SearchByISBN(
	ctx context.Context,
	req *searcherv1.SearchByISBNRequest,
) (*searcherv1.SearchByISBNResponse, error) {
	return &searcherv1.SearchByISBNResponse{
		Title:     "title",
		Author:    "author",
		Publisher: "publisher",
		Year:      "year",
	}, nil
}
