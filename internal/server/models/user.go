package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

// User represents a user in the system.
type User struct {
	ID          string    `bson:"_id" json:"id"` // Firebase UID as primary key
	DisplayName string    `bson:"display_name" json:"display_name"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}

func (u *User) ToMap() bson.M {
	return bson.M{
		"display_name": u.DisplayName,
		"created_at":   u.CreatedAt,
		"updated_at":   u.UpdatedAt,
	}
}
