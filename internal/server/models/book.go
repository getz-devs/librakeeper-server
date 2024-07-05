package models

import (
	"time"
)

// Book represents a book in the library.
type Book struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	UserID      string    `bson:"user_id" json:"user_id"`
	BookshelfID string    `bson:"bookshelf_id" json:"bookshelf_id"`
	Title       string    `bson:"title" json:"title"`
	Author      string    `bson:"author" json:"author"`
	ISBN        string    `bson:"isbn" json:"isbn"`
	Description string    `bson:"description" json:"description"`
	CoverImage  string    `bson:"cover_image" json:"cover_image"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}

// BookUpdate represents fields that can be updated in a Book.
type BookUpdate struct {
	Title       *string   `bson:"title,omitempty" json:"title,omitempty"` // Optional fields for update
	Author      *string   `bson:"author,omitempty" json:"author,omitempty"`
	ISBN        *string   `bson:"isbn,omitempty" json:"isbn,omitempty"`
	Description *string   `bson:"description,omitempty" json:"description,omitempty"`
	CoverImage  *string   `bson:"cover_image,omitempty" json:"cover_image,omitempty"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}
