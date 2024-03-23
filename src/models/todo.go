package models

import (
    "time"
    "github.com/google/uuid"
)

type Todo struct {
    ID          uuid.UUID    `json:"id"` 
    Title       string       `json:"title"`
    Completed   bool         `json:"completed"` // Default: false
    CreatedAt   time.Time    `json:"created_at,omitempty"` // Default: time.Now()
    UpdatedAt   time.Time    `json:"updated_at,omitempty"` // Default: time.Now() 
}
