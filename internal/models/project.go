package models

import (
	"time"
)

// Project represents a project in the system
type Project struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OwnerID     int64     `json:"owner_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ProjectCreate represents data needed to create a new project
type ProjectCreate struct {
	Name        string `json:"name" binding:"required,min=3,max=100"`
	Description string `json:"description" binding:"max=500"`
}

// ProjectUpdate represents data needed to update a project
type ProjectUpdate struct {
	Name        string `json:"name" binding:"omitempty,min=3,max=100"`
	Description string `json:"description" binding:"max=500"`
}

// ProjectResponse represents the project data to be returned in API responses
type ProjectResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OwnerID     int64     `json:"owner_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToResponse converts a Project to ProjectResponse
func (p *Project) ToResponse() ProjectResponse {
	return ProjectResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		OwnerID:     p.OwnerID,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
