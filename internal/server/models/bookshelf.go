package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Bookshelf represents a collection of books.
type Bookshelf struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    string             `bson:"user_id" json:"user_id"`
	Name      string             `bson:"name" json:"name"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
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
