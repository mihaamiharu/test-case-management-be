package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/mihaamiharu/test-case-management-be/internal/models"
)

var (
	ErrTagNotFound = errors.New("tag not found")
	ErrTagExists   = errors.New("tag with this name already exists")
)

// TagRepositoryInterface defines the interface for tag repository operations
type TagRepositoryInterface interface {
	Create(tag *models.Tag) error
	GetByID(id int64) (*models.Tag, error)
	GetByName(name string) (*models.Tag, error)
	Delete(id int64) error
	List() ([]*models.Tag, error)
	GetTagsByTestCase(testCaseID int64) ([]*models.Tag, error)
	AddTagToTestCase(testCaseID, tagID int64) error
	RemoveTagFromTestCase(testCaseID, tagID int64) error
	UpdateTestCaseTags(testCaseID int64, tagIDs []int64) error
	GetOrCreateTag(name string) (*models.Tag, error)
}

// TagRepository handles database operations for tags
type TagRepository struct {
	db *sql.DB
}

// NewTagRepository creates a new tag repository
func NewTagRepository(db *sql.DB) *TagRepository {
	return &TagRepository{db: db}
}

// Create adds a new tag to the database
func (r *TagRepository) Create(tag *models.Tag) error {
	// Check if tag with name already exists
	existingTag, err := r.GetByName(tag.Name)
	if err == nil && existingTag != nil {
		return ErrTagExists
	}
	if err != ErrTagNotFound {
		return fmt.Errorf("failed to check existing tag: %v", err)
	}

	query := `
		INSERT INTO tags (name, created_at)
		VALUES (?, ?)`

	now := time.Now()
	result, err := r.db.Exec(query, tag.Name, now)
	if err != nil {
		return fmt.Errorf("failed to create tag: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %v", err)
	}

	tag.ID = id
	tag.CreatedAt = now
	return nil
}

// GetByID retrieves a tag by ID
func (r *TagRepository) GetByID(id int64) (*models.Tag, error) {
	tag := &models.Tag{}
	query := `
		SELECT id, name, created_at
		FROM tags
		WHERE id = ?`

	err := r.db.QueryRow(query, id).Scan(&tag.ID, &tag.Name, &tag.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrTagNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get tag: %v", err)
	}

	return tag, nil
}

// GetByName retrieves a tag by name
func (r *TagRepository) GetByName(name string) (*models.Tag, error) {
	tag := &models.Tag{}
	query := `
		SELECT id, name, created_at
		FROM tags
		WHERE name = ?`

	err := r.db.QueryRow(query, name).Scan(&tag.ID, &tag.Name, &tag.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrTagNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get tag: %v", err)
	}

	return tag, nil
}

// Delete removes a tag from the database
func (r *TagRepository) Delete(id int64) error {
	result, err := r.db.Exec("DELETE FROM tags WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete tag: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return ErrTagNotFound
	}

	return nil
}

// List retrieves all tags
func (r *TagRepository) List() ([]*models.Tag, error) {
	query := `
		SELECT id, name, created_at
		FROM tags
		ORDER BY name`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %v", err)
	}
	defer rows.Close()

	var tags []*models.Tag
	for rows.Next() {
		tag := &models.Tag{}
		err := rows.Scan(&tag.ID, &tag.Name, &tag.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tag: %v", err)
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

// GetTagsByTestCase retrieves all tags for a specific test case
func (r *TagRepository) GetTagsByTestCase(testCaseID int64) ([]*models.Tag, error) {
	query := `
		SELECT t.id, t.name, t.created_at
		FROM tags t
		JOIN test_case_tags tct ON t.id = tct.tag_id
		WHERE tct.test_case_id = ?
		ORDER BY t.name`

	rows, err := r.db.Query(query, testCaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags for test case: %v", err)
	}
	defer rows.Close()

	var tags []*models.Tag
	for rows.Next() {
		tag := &models.Tag{}
		err := rows.Scan(&tag.ID, &tag.Name, &tag.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tag: %v", err)
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

// AddTagToTestCase adds a tag to a test case
func (r *TagRepository) AddTagToTestCase(testCaseID, tagID int64) error {
	// Check if the association already exists
	var count int
	err := r.db.QueryRow(
		"SELECT COUNT(*) FROM test_case_tags WHERE test_case_id = ? AND tag_id = ?",
		testCaseID, tagID,
	).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check existing tag association: %v", err)
	}

	if count > 0 {
		// Association already exists, no need to insert
		return nil
	}

	query := `
		INSERT INTO test_case_tags (test_case_id, tag_id)
		VALUES (?, ?)`

	_, err = r.db.Exec(query, testCaseID, tagID)
	if err != nil {
		return fmt.Errorf("failed to add tag to test case: %v", err)
	}

	return nil
}

// RemoveTagFromTestCase removes a tag from a test case
func (r *TagRepository) RemoveTagFromTestCase(testCaseID, tagID int64) error {
	result, err := r.db.Exec(
		"DELETE FROM test_case_tags WHERE test_case_id = ? AND tag_id = ?",
		testCaseID, tagID,
	)
	if err != nil {
		return fmt.Errorf("failed to remove tag from test case: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("tag was not associated with test case")
	}

	return nil
}

// UpdateTestCaseTags updates the tags for a test case
func (r *TagRepository) UpdateTestCaseTags(testCaseID int64, tagIDs []int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Remove all existing tags
	_, err = tx.Exec("DELETE FROM test_case_tags WHERE test_case_id = ?", testCaseID)
	if err != nil {
		return fmt.Errorf("failed to remove existing tags: %v", err)
	}

	// Add new tags
	if len(tagIDs) > 0 {
		query := `
			INSERT INTO test_case_tags (test_case_id, tag_id)
			VALUES (?, ?)`

		for _, tagID := range tagIDs {
			_, err = tx.Exec(query, testCaseID, tagID)
			if err != nil {
				return fmt.Errorf("failed to add tag %d: %v", tagID, err)
			}
		}
	}

	return tx.Commit()
}

func (r *TagRepository) GetOrCreateTag(name string) (*models.Tag, error) {
	// First try to get the existing tag
	tag, err := r.GetByName(name)
	if err == nil {
		return tag, nil
	}
	if err != ErrTagNotFound {
		return nil, err
	}

	// Tag doesn't exist, create it
	tag = &models.Tag{Name: name}
	err = r.Create(tag)
	if err != nil {
		return nil, err
	}

	return tag, nil
}
