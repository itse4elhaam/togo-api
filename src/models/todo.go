package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Todo struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Title     string             `json:"title"`
	Completed bool               `json:"completed"`
	CreatedAt time.Time          `json:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at,omitempty"`
}
