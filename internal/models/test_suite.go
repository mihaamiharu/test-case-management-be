package models

import (
	"time"
)

// TestSuite represents a collection of test cases
type TestSuite struct {
	ID          int64     `json:"id"`
	ProjectID   int64     `json:"project_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TestSuiteCreate represents data needed to create a new test suite
type TestSuiteCreate struct {
	ProjectID   int64  `json:"project_id" binding:"required"`
	Name        string `json:"name" binding:"required,min=3,max=100"`
	Description string `json:"description" binding:"max=1000"`
}

// TestSuiteUpdate represents data needed to update a test suite
type TestSuiteUpdate struct {
	Name        string `json:"name" binding:"omitempty,min=3,max=100"`
	Description string `json:"description" binding:"max=1000"`
}

// TestSuiteResponse represents the test suite data to be returned in API responses
type TestSuiteResponse struct {
	ID          int64     `json:"id"`
	ProjectID   int64     `json:"project_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToResponse converts a TestSuite to TestSuiteResponse
func (ts *TestSuite) ToResponse() TestSuiteResponse {
	return TestSuiteResponse{
		ID:          ts.ID,
		ProjectID:   ts.ProjectID,
		Name:        ts.Name,
		Description: ts.Description,
		CreatedAt:   ts.CreatedAt,
		UpdatedAt:   ts.UpdatedAt,
	}
}
