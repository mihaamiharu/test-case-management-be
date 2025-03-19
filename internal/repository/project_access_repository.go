package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/mihaamiharu/test-case-management-be/internal/models"
)

var (
	ErrProjectAccessNotFound = errors.New("project access not found")
	ErrProjectAccessExists   = errors.New("user already has access to this project")
)

// ProjectAccessRepositoryInterface defines the interface for project access repository operations
type ProjectAccessRepositoryInterface interface {
	Create(access *models.ProjectAccess) error
	GetByID(id int64) (*models.ProjectAccess, error)
	GetByProjectAndUser(projectID, userID int64) (*models.ProjectAccess, error)
	Update(access *models.ProjectAccess) error
	Delete(id int64) error
	ListByProject(projectID int64) ([]*models.ProjectAccess, error)
	ListByUser(userID int64) ([]*models.ProjectAccess, error)
	GetProjectIDsByUserAccess(userID int64) ([]int64, error)
	HasEditAccess(projectID, userID int64) (bool, error)
	HasViewAccess(projectID, userID int64) (bool, error)
}

// ProjectAccessRepository handles database operations for project access
type ProjectAccessRepository struct {
	db *sql.DB
}

// NewProjectAccessRepository creates a new project access repository
func NewProjectAccessRepository(db *sql.DB) *ProjectAccessRepository {
	return &ProjectAccessRepository{db: db}
}

// Create adds a new project access record to the database
func (r *ProjectAccessRepository) Create(access *models.ProjectAccess) error {
	// Check if access already exists
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM project_access WHERE project_id = ? AND user_id = ?",
		access.ProjectID, access.UserID).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrProjectAccessExists
	}

	// Insert new access record
	query := `
		INSERT INTO project_access (project_id, user_id, level, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`
	now := time.Now()
	result, err := r.db.Exec(query, access.ProjectID, access.UserID, access.Level, now, now)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	access.ID = id
	access.CreatedAt = now
	access.UpdatedAt = now
	return nil
}

// GetByID retrieves a project access record by ID
func (r *ProjectAccessRepository) GetByID(id int64) (*models.ProjectAccess, error) {
	query := `
		SELECT id, project_id, user_id, level, created_at, updated_at
		FROM project_access
		WHERE id = ?
	`
	access := &models.ProjectAccess{}
	err := r.db.QueryRow(query, id).Scan(
		&access.ID,
		&access.ProjectID,
		&access.UserID,
		&access.Level,
		&access.CreatedAt,
		&access.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrProjectAccessNotFound
		}
		return nil, err
	}
	return access, nil
}

// GetByProjectAndUser retrieves a project access record by project ID and user ID
func (r *ProjectAccessRepository) GetByProjectAndUser(projectID, userID int64) (*models.ProjectAccess, error) {
	query := `
		SELECT id, project_id, user_id, level, created_at, updated_at
		FROM project_access
		WHERE project_id = ? AND user_id = ?
	`
	access := &models.ProjectAccess{}
	err := r.db.QueryRow(query, projectID, userID).Scan(
		&access.ID,
		&access.ProjectID,
		&access.UserID,
		&access.Level,
		&access.CreatedAt,
		&access.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrProjectAccessNotFound
		}
		return nil, err
	}
	return access, nil
}

// Update updates an existing project access record
func (r *ProjectAccessRepository) Update(access *models.ProjectAccess) error {
	// Check if access exists
	_, err := r.GetByID(access.ID)
	if err != nil {
		return err
	}

	// Update access
	query := `
		UPDATE project_access
		SET level = ?, updated_at = ?
		WHERE id = ?
	`
	now := time.Now()
	_, err = r.db.Exec(query, access.Level, now, access.ID)
	if err != nil {
		return err
	}

	access.UpdatedAt = now
	return nil
}

// Delete removes a project access record from the database
func (r *ProjectAccessRepository) Delete(id int64) error {
	// Check if access exists
	_, err := r.GetByID(id)
	if err != nil {
		return err
	}

	// Delete access
	query := `DELETE FROM project_access WHERE id = ?`
	_, err = r.db.Exec(query, id)
	return err
}

// ListByProject retrieves all access records for a specific project
func (r *ProjectAccessRepository) ListByProject(projectID int64) ([]*models.ProjectAccess, error) {
	query := `
		SELECT id, project_id, user_id, level, created_at, updated_at
		FROM project_access
		WHERE project_id = ?
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accessList []*models.ProjectAccess
	for rows.Next() {
		access := &models.ProjectAccess{}
		err := rows.Scan(
			&access.ID,
			&access.ProjectID,
			&access.UserID,
			&access.Level,
			&access.CreatedAt,
			&access.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		accessList = append(accessList, access)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return accessList, nil
}

// ListByUser retrieves all access records for a specific user
func (r *ProjectAccessRepository) ListByUser(userID int64) ([]*models.ProjectAccess, error) {
	query := `
		SELECT id, project_id, user_id, level, created_at, updated_at
		FROM project_access
		WHERE user_id = ?
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accessList []*models.ProjectAccess
	for rows.Next() {
		access := &models.ProjectAccess{}
		err := rows.Scan(
			&access.ID,
			&access.ProjectID,
			&access.UserID,
			&access.Level,
			&access.CreatedAt,
			&access.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		accessList = append(accessList, access)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return accessList, nil
}

// GetProjectIDsByUserAccess retrieves all project IDs that a user has access to
func (r *ProjectAccessRepository) GetProjectIDsByUserAccess(userID int64) ([]int64, error) {
	query := `
		SELECT project_id
		FROM project_access
		WHERE user_id = ?
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projectIDs []int64
	for rows.Next() {
		var projectID int64
		err := rows.Scan(&projectID)
		if err != nil {
			return nil, err
		}
		projectIDs = append(projectIDs, projectID)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return projectIDs, nil
}

// HasEditAccess checks if a user has edit access to a project
func (r *ProjectAccessRepository) HasEditAccess(projectID, userID int64) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM project_access
		WHERE project_id = ? AND user_id = ? AND level = 'edit'
	`
	var count int
	err := r.db.QueryRow(query, projectID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// HasViewAccess checks if a user has at least view access to a project
func (r *ProjectAccessRepository) HasViewAccess(projectID, userID int64) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM project_access
		WHERE project_id = ? AND user_id = ?
	`
	var count int
	err := r.db.QueryRow(query, projectID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
