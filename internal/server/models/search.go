package models

import searcherv1 "github.com/getz-devs/librakeeper-protos/gen/go/searcher"

type SearchResponse struct {
	Status searcherv1.SearchByISBNResponse_Status `json:"status"`
	Books  []*Book                                `json:"books"`
}
