package models

import (
	"time"
)

// Book represents a book in the library.
type Book struct {
	ID          string `bson:"_id,omitempty" json:"id"`
	UserID      string `bson:"user_id" json:"user_id"`
	BookshelfID string `bson:"bookshelf_id" json:"bookshelf_id"`

	ISBN        string `bson:"isbn" json:"isbn"`
	Title       string `bson:"title" json:"title"`
	Author      string `bson:"author" json:"author"`
	Publishing  string `bson:"publishing" json:"publishing"`
	Description string `bson:"description" json:"description"`
	CoverImage  string `bson:"cover_image" json:"cover_image"`
	ShopName    string `bson:"shop_name" json:"shop_name"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// BookUpdate represents fields that can be updated in a Book.
type BookUpdate struct {
	ISBN        *string   `bson:"isbn,omitempty" json:"isbn,omitempty"`
	Title       *string   `bson:"title,omitempty" json:"title,omitempty"` // Optional fields for update
	Author      *string   `bson:"author,omitempty" json:"author,omitempty"`
	Publishing  *string   `bson:"publishing" json:"publishing"`
	Description *string   `bson:"description,omitempty" json:"description,omitempty"`
	CoverImage  *string   `bson:"cover_image,omitempty" json:"cover_image,omitempty"`
	ShopName    *string   `bson:"shop_name" json:"shop_name"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}
