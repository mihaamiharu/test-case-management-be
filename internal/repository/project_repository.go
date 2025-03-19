package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/mihaamiharu/test-case-management-be/internal/models"
)

var (
	ErrProjectNotFound = errors.New("project not found")
	ErrProjectExists   = errors.New("project with this name already exists")
)

// ProjectRepositoryInterface defines the interface for project repository operations
type ProjectRepositoryInterface interface {
	Create(project *models.Project) error
	GetByID(id int64) (*models.Project, error)
	Update(project *models.Project) error
	Delete(id int64) error
	ListByOwner(ownerID int64) ([]*models.Project, error)
	List(page, pageSize int) ([]*models.Project, error)
	IsOwner(projectID, userID int64) (bool, error)
}

// ProjectRepository handles database operations for projects
type ProjectRepository struct {
	db *sql.DB
}

// NewProjectRepository creates a new project repository
func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

// Create adds a new project to the database
func (r *ProjectRepository) Create(project *models.Project) error {
	// Check if project with name already exists for this owner
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM projects WHERE name = ? AND owner_id = ?",
		project.Name, project.OwnerID).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrProjectExists
	}

	// Insert new project
	query := `
		INSERT INTO projects (name, description, owner_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`
	now := time.Now()
	result, err := r.db.Exec(query, project.Name, project.Description, project.OwnerID, now, now)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	project.ID = id
	project.CreatedAt = now
	project.UpdatedAt = now
	return nil
}

// GetByID retrieves a project by ID
func (r *ProjectRepository) GetByID(id int64) (*models.Project, error) {
	query := `
		SELECT id, name, description, owner_id, created_at, updated_at
		FROM projects
		WHERE id = ?
	`
	project := &models.Project{}
	err := r.db.QueryRow(query, id).Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.OwnerID,
		&project.CreatedAt,
		&project.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}
	return project, nil
}

// Update updates an existing project
func (r *ProjectRepository) Update(project *models.Project) error {
	// Check if project exists
	_, err := r.GetByID(project.ID)
	if err != nil {
		return err
	}

	// Update project
	query := `
		UPDATE projects
		SET name = ?, description = ?, updated_at = ?
		WHERE id = ?
	`
	now := time.Now()
	_, err = r.db.Exec(query, project.Name, project.Description, now, project.ID)
	if err != nil {
		return err
	}

	project.UpdatedAt = now
	return nil
}

// Delete removes a project from the database
func (r *ProjectRepository) Delete(id int64) error {
	// Check if project exists
	_, err := r.GetByID(id)
	if err != nil {
		return err
	}

	// Delete project
	query := `DELETE FROM projects WHERE id = ?`
	_, err = r.db.Exec(query, id)
	return err
}

// ListByOwner retrieves all projects for a specific owner
func (r *ProjectRepository) ListByOwner(ownerID int64) ([]*models.Project, error) {
	query := `
		SELECT id, name, description, owner_id, created_at, updated_at
		FROM projects
		WHERE owner_id = ?
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*models.Project
	for rows.Next() {
		project := &models.Project{}
		err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.Description,
			&project.OwnerID,
			&project.CreatedAt,
			&project.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

// List retrieves all projects with optional pagination
func (r *ProjectRepository) List(limit, offset int) ([]*models.Project, error) {
	query := `
		SELECT id, name, description, owner_id, created_at, updated_at
		FROM projects
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*models.Project
	for rows.Next() {
		project := &models.Project{}
		err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.Description,
			&project.OwnerID,
			&project.CreatedAt,
			&project.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

// IsOwner checks if a user is the owner of a project
func (r *ProjectRepository) IsOwner(projectID, userID int64) (bool, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM projects WHERE id = ? AND owner_id = ?", projectID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
