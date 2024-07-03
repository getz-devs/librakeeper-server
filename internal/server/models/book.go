package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Book represents a book in the library.
type Book struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BookshelfID primitive.ObjectID `bson:"bookshelf_id" json:"bookshelf_id"`
	Title       string             `bson:"title" json:"title"`
	Author      string             `bson:"author" json:"author"`
	ISBN        string             `bson:"isbn" json:"isbn"`
	Description string             `bson:"description" json:"description"`
	CoverImage  string             `bson:"cover_image" json:"cover_image"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// ToMap converts the Book struct to a bson.M map for MongoDB updates.
func (b *Book) ToMap() bson.M {
	return bson.M{
		"bookshelf_id": b.BookshelfID,
		"title":        b.Title,
		"author":       b.Author,
		"isbn":         b.ISBN,
		"description":  b.Description,
		"cover_image":  b.CoverImage,
		"created_at":   b.CreatedAt,
		"updated_at":   b.UpdatedAt,
	}
}
