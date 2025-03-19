package models

import (
	"time"
)

// AccessLevel defines the level of access a user has to a project
type AccessLevel string

const (
	// AccessLevelView allows a user to view a project but not modify it
	AccessLevelView AccessLevel = "view"

	// AccessLevelEdit allows a user to view and edit a project
	AccessLevelEdit AccessLevel = "edit"
)

// ProjectAccess represents a user's access to a project
type ProjectAccess struct {
	ID        int64       `json:"id"`
	ProjectID int64       `json:"project_id"`
	UserID    int64       `json:"user_id"`
	Level     AccessLevel `json:"level"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// ProjectAccessCreate represents data needed to grant access to a project
type ProjectAccessCreate struct {
	UserID int64       `json:"user_id" binding:"required"`
	Level  AccessLevel `json:"level" binding:"required,oneof=view edit"`
}

// ProjectAccessUpdate represents data needed to update project access
type ProjectAccessUpdate struct {
	Level AccessLevel `json:"level" binding:"required,oneof=view edit"`
}

// ProjectAccessResponse represents the project access data to be returned in API responses
type ProjectAccessResponse struct {
	ID        int64       `json:"id"`
	ProjectID int64       `json:"project_id"`
	UserID    int64       `json:"user_id"`
	Level     AccessLevel `json:"level"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// ToResponse converts a ProjectAccess to ProjectAccessResponse
func (pa *ProjectAccess) ToResponse() ProjectAccessResponse {
	return ProjectAccessResponse{
		ID:        pa.ID,
		ProjectID: pa.ProjectID,
		UserID:    pa.UserID,
		Level:     pa.Level,
		CreatedAt: pa.CreatedAt,
		UpdatedAt: pa.UpdatedAt,
	}
}
