package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/mihaamiharu/test-case-management-be/internal/models"
)

var (
	ErrTestSuiteNotFound = errors.New("test suite not found")
	ErrTestSuiteExists   = errors.New("test suite with this name already exists in the project")
)

// TestSuiteRepositoryInterface defines the interface for test suite repository operations
type TestSuiteRepositoryInterface interface {
	Create(suite *models.TestSuite) error
	GetByID(id int64) (*models.TestSuite, error)
	Update(suite *models.TestSuite) error
	Delete(id int64) error
	ListByProject(projectID int64) ([]*models.TestSuite, error)
	List() ([]*models.TestSuite, error)
}

// TestSuiteRepository handles database operations for test suites
type TestSuiteRepository struct {
	db *sql.DB
}

// NewTestSuiteRepository creates a new test suite repository
func NewTestSuiteRepository(db *sql.DB) *TestSuiteRepository {
	return &TestSuiteRepository{db: db}
}

// Create adds a new test suite to the database
func (r *TestSuiteRepository) Create(suite *models.TestSuite) error {
	// Check if suite with name already exists in this project
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM test_suites WHERE name = ? AND project_id = ?",
		suite.Name, suite.ProjectID).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrTestSuiteExists
	}

	// Insert new test suite
	query := `
		INSERT INTO test_suites (project_id, name, description, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`
	now := time.Now()
	result, err := r.db.Exec(query, suite.ProjectID, suite.Name, suite.Description, now, now)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	suite.ID = id
	suite.CreatedAt = now
	suite.UpdatedAt = now
	return nil
}

// GetByID retrieves a test suite by ID
func (r *TestSuiteRepository) GetByID(id int64) (*models.TestSuite, error) {
	query := `
		SELECT id, project_id, name, description, created_at, updated_at
		FROM test_suites
		WHERE id = ?
	`
	suite := &models.TestSuite{}
	err := r.db.QueryRow(query, id).Scan(
		&suite.ID,
		&suite.ProjectID,
		&suite.Name,
		&suite.Description,
		&suite.CreatedAt,
		&suite.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrTestSuiteNotFound
		}
		return nil, err
	}
	return suite, nil
}

// Update updates an existing test suite
func (r *TestSuiteRepository) Update(suite *models.TestSuite) error {
	// Check if suite exists
	_, err := r.GetByID(suite.ID)
	if err != nil {
		return err
	}

	// Check if the new name conflicts with another suite in the same project
	if suite.Name != "" {
		var count int
		err := r.db.QueryRow("SELECT COUNT(*) FROM test_suites WHERE name = ? AND project_id = ? AND id != ?",
			suite.Name, suite.ProjectID, suite.ID).Scan(&count)
		if err != nil {
			return err
		}
		if count > 0 {
			return ErrTestSuiteExists
		}
	}

	// Update suite
	query := `
		UPDATE test_suites
		SET name = ?, description = ?, updated_at = ?
		WHERE id = ?
	`
	now := time.Now()
	_, err = r.db.Exec(query, suite.Name, suite.Description, now, suite.ID)
	if err != nil {
		return err
	}

	suite.UpdatedAt = now
	return nil
}

// Delete removes a test suite from the database
func (r *TestSuiteRepository) Delete(id int64) error {
	// Check if suite exists
	_, err := r.GetByID(id)
	if err != nil {
		return err
	}

	// Delete suite
	query := `DELETE FROM test_suites WHERE id = ?`
	_, err = r.db.Exec(query, id)
	return err
}

// ListByProject retrieves all test suites for a specific project
func (r *TestSuiteRepository) ListByProject(projectID int64) ([]*models.TestSuite, error) {
	query := `
		SELECT id, project_id, name, description, created_at, updated_at
		FROM test_suites
		WHERE project_id = ?
		ORDER BY name ASC
	`
	rows, err := r.db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var suites []*models.TestSuite
	for rows.Next() {
		suite := &models.TestSuite{}
		err := rows.Scan(
			&suite.ID,
			&suite.ProjectID,
			&suite.Name,
			&suite.Description,
			&suite.CreatedAt,
			&suite.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		suites = append(suites, suite)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return suites, nil
}

// List retrieves all test suites
func (r *TestSuiteRepository) List() ([]*models.TestSuite, error) {
	query := `
		SELECT id, project_id, name, description, created_at, updated_at
		FROM test_suites
		ORDER BY name ASC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var suites []*models.TestSuite
	for rows.Next() {
		suite := &models.TestSuite{}
		err := rows.Scan(
			&suite.ID,
			&suite.ProjectID,
			&suite.Name,
			&suite.Description,
			&suite.CreatedAt,
			&suite.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		suites = append(suites, suite)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return suites, nil
}
