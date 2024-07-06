package models

import (
	"time"
)

// User represents a user in the system.
type User struct {
	ID          string    `bson:"_id,omitempty" json:"id"` // Firebase UID as primary key
	DisplayName string    `bson:"display_name" json:"display_name"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}

// UserUpdate represents fields that can be updated in a User.
type UserUpdate struct {
	DisplayName *string   `bson:"display_name,omitempty" json:"display_name,omitempty"` // Optional field for update
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}
