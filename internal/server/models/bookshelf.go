package models

import (
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

// BookshelfUpdate represents fields that can be updated in a Bookshelf.
type BookshelfUpdate struct {
	Name      *string   `bson:"name,omitempty" json:"name,omitempty"` // Optional field for update
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
