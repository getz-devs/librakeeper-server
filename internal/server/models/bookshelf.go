package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

// Bookshelf represents a collection of books.
type Bookshelf struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	UserID    string    `bson:"user_id" json:"user_id"`
	Name      string    `bson:"name" json:"name"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// ToMap converts the Bookshelf struct to a bson.M map for MongoDB updates.
func (b *Bookshelf) ToMap() bson.M {
	return bson.M{
		"user_id":    b.UserID,
		"name":       b.Name,
		"created_at": b.CreatedAt,
		"updated_at": b.UpdatedAt,
	}
}

// BookshelfUpdate represents fields that can be updated in a Bookshelf.
type BookshelfUpdate struct {
	Name *string `bson:"name,omitempty" json:"name,omitempty"` // Optional field for update
}
