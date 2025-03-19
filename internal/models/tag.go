package models

import (
	"time"
)

// Tag represents a label that can be applied to test cases
type Tag struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// TagCreate represents data needed to create a new tag
type TagCreate struct {
	Name string `json:"name" binding:"required,min=1,max=50"`
}

// TagResponse represents the tag data to be returned in API responses
type TagResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse converts a Tag to TagResponse
func (t *Tag) ToResponse() *TagResponse {
	return &TagResponse{
		ID:        t.ID,
		Name:      t.Name,
		CreatedAt: t.CreatedAt,
	}
}
