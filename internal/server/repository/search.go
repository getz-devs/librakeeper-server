package repository

import (
	"context"
	searcherv1 "github.com/getz-devs/librakeeper-protos/gen/go/searcher"
)

type SearchRepo interface {
	SearchByISBN(ctx context.Context, isbn string) (*searcherv1.SearchByISBNResponse, error)
}
