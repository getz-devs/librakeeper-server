package bookModels

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type BookInShop struct {
	Title      string `selector:"div.results__book-name > a" bson:"title,omitempty"`
	Author     string `selector:"div.results__authors" bson:"author,omitempty"`
	Publishing string `selector:"div.results__publishing" bson:"publishing,omitempty"`
	ImgUrl     string `selector:"a.results__image > img" attr:"src" bson:"img_url"`
	ShopName   string `selector:"div.results__shop-name > a" bson:"shop_name"`
}

type RequestStatus int

// Pending, Success, Failed
const (
	Pending RequestStatus = iota
	Success
	Failed
)

type SearchRequest struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Isbn      string             `bson:"isbn"`
	Status    RequestStatus      `bson:"status,omitempty"`
	Books     []*BookInShop      `bson:"books,omitempty"`
	CreatedAt primitive.DateTime `bson:"created_at,omitempty"`
	UpdatedAt primitive.DateTime `bson:"updated_at,omitempty"`
}

// New creates a new SearchRequest with the provided ISBN, current time, and initial values.
//
// Parameters:
//
//	isbn - the ISBN of the book
//
// Returns:
//
//	SearchRequest - the newly created SearchRequest
func New(isbn string) SearchRequest {
	currentTime := primitive.NewDateTimeFromTime(time.Now())
	return SearchRequest{
		ID:        primitive.NewObjectID(),
		Isbn:      isbn,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Status:    Pending,
		Books:     []*BookInShop{},
	}
}
