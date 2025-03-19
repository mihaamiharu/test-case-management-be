package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/mihaamiharu/test-case-management-be/internal/models"
)

var (
	ErrTestCaseNotFound = errors.New("test case not found")
)

type TestCaseRepositoryInterface interface {
	Create(testCase *models.TestCase) error
	GetByID(id int64) (*models.TestCase, error)
	Update(testCase *models.TestCase) error
	Delete(id int64) error
	ListByProject(projectID int64) ([]*models.TestCase, error)
	ListBySuite(suiteID int64) ([]*models.TestCase, error)
	GetSteps(testCaseID int64) ([]*models.TestStep, error)
	CreateStep(step *models.TestStep) error
	UpdateStep(step *models.TestStep) error
	DeleteStep(stepID int64) error
	CreateStepNote(note *models.StepNote) error
	DeleteStepNote(noteID int64) error
	CreateStepAttachment(attachment *models.StepAttachment) error
	DeleteStepAttachment(attachmentID int64) error
	GetStepByID(stepID int64) (*models.TestStep, error)
	GetStepAttachmentByID(attachmentID int64) (*models.StepAttachment, error)
}

type TestCaseRepository struct {
	db *sql.DB
}

func NewTestCaseRepository(db *sql.DB) *TestCaseRepository {
	return &TestCaseRepository{db: db}
}

func (r *TestCaseRepository) Create(testCase *models.TestCase) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	now := time.Now()
	// Insert test case
	query := `
		INSERT INTO test_cases (
			project_id, suite_id, title, description, preconditions,
			status, priority, created_by, updated_by, version,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := tx.Exec(
		query,
		testCase.ProjectID,
		testCase.SuiteID,
		testCase.Title,
		testCase.Description,
		testCase.Preconditions,
		testCase.Status,
		testCase.Priority,
		testCase.CreatedBy,
		testCase.UpdatedBy,
		1, // Initial version
		now,
		now,
	)

	if err != nil {
		if err.Error() == "Error 1062: Duplicate entry" {
			return fmt.Errorf("test case with title '%s' already exists in this project", testCase.Title)
		}
		return fmt.Errorf("failed to create test case: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %v", err)
	}

	testCase.ID = id
	testCase.CreatedAt = now
	testCase.UpdatedAt = now

	// Insert steps if any
	if len(testCase.Steps) > 0 {
		stepQuery := `
			INSERT INTO test_steps (
				test_case_id, step_number, step_type, description,
				expected_result, created_at, updated_at
			) VALUES (?, ?, ?, ?, ?, ?, ?)`

		for i, step := range testCase.Steps {
			step.TestCaseID = testCase.ID
			step.StepNumber = i + 1

			stepResult, err := tx.Exec(
				stepQuery,
				step.TestCaseID,
				step.StepNumber,
				step.StepType,
				step.Description,
				step.ExpectedResult,
				now,
				now,
			)

			if err != nil {
				return fmt.Errorf("failed to create test step: %v", err)
			}

			stepID, err := stepResult.LastInsertId()
			if err != nil {
				return fmt.Errorf("failed to get step last insert ID: %v", err)
			}

			step.ID = stepID
			step.CreatedAt = now
			step.UpdatedAt = now
		}
	}

	return tx.Commit()
}

func (r *TestCaseRepository) GetByID(id int64) (*models.TestCase, error) {
	testCase := &models.TestCase{}
	query := `
		SELECT 
			id, project_id, suite_id, title, description, preconditions,
			status, priority, created_by, updated_by, version,
			created_at, updated_at
		FROM test_cases
		WHERE id = ?`

	err := r.db.QueryRow(query, id).Scan(
		&testCase.ID,
		&testCase.ProjectID,
		&testCase.SuiteID,
		&testCase.Title,
		&testCase.Description,
		&testCase.Preconditions,
		&testCase.Status,
		&testCase.Priority,
		&testCase.CreatedBy,
		&testCase.UpdatedBy,
		&testCase.Version,
		&testCase.CreatedAt,
		&testCase.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrTestCaseNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get test case: %v", err)
	}

	// Get steps
	steps, err := r.GetSteps(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get test steps: %v", err)
	}
	testCase.Steps = steps

	return testCase, nil
}

func (r *TestCaseRepository) Update(testCase *models.TestCase) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	now := time.Now()
	// Update test case
	query := `
		UPDATE test_cases SET
			project_id = ?,
			suite_id = ?,
			title = ?,
			description = ?,
			preconditions = ?,
			status = ?,
			priority = ?,
			updated_by = ?,
			version = version + 1,
			updated_at = ?
		WHERE id = ?`

	result, err := tx.Exec(
		query,
		testCase.ProjectID,
		testCase.SuiteID,
		testCase.Title,
		testCase.Description,
		testCase.Preconditions,
		testCase.Status,
		testCase.Priority,
		testCase.UpdatedBy,
		now,
		testCase.ID,
	)

	if err != nil {
		if err.Error() == "Error 1062: Duplicate entry" {
			return fmt.Errorf("test case with title '%s' already exists in this project", testCase.Title)
		}
		return fmt.Errorf("failed to update test case: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return ErrTestCaseNotFound
	}

	testCase.UpdatedAt = now
	testCase.Version++

	// Delete existing steps
	_, err = tx.Exec("DELETE FROM test_steps WHERE test_case_id = ?", testCase.ID)
	if err != nil {
		return fmt.Errorf("failed to delete existing steps: %v", err)
	}

	// Insert updated steps
	if len(testCase.Steps) > 0 {
		stepQuery := `
			INSERT INTO test_steps (
				test_case_id, step_number, step_type, description,
				expected_result, created_at, updated_at
			) VALUES (?, ?, ?, ?, ?, ?, ?)`

		for i, step := range testCase.Steps {
			step.TestCaseID = testCase.ID
			step.StepNumber = i + 1

			stepResult, err := tx.Exec(
				stepQuery,
				step.TestCaseID,
				step.StepNumber,
				step.StepType,
				step.Description,
				step.ExpectedResult,
				now,
				now,
			)

			if err != nil {
				return fmt.Errorf("failed to create test step: %v", err)
			}

			stepID, err := stepResult.LastInsertId()
			if err != nil {
				return fmt.Errorf("failed to get step last insert ID: %v", err)
			}

			step.ID = stepID
			step.CreatedAt = now
			step.UpdatedAt = now
		}
	}

	return tx.Commit()
}

func (r *TestCaseRepository) Delete(id int64) error {
	result, err := r.db.Exec("DELETE FROM test_cases WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete test case: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return ErrTestCaseNotFound
	}

	return nil
}

func (r *TestCaseRepository) ListByProject(projectID int64) ([]*models.TestCase, error) {
	query := `
		SELECT 
			id, project_id, suite_id, title, description, preconditions,
			status, priority, created_by, updated_by, version,
			created_at, updated_at
		FROM test_cases
		WHERE project_id = ?
		ORDER BY title`

	rows, err := r.db.Query(query, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list test cases: %v", err)
	}
	defer rows.Close()

	var testCases []*models.TestCase
	for rows.Next() {
		testCase := &models.TestCase{}
		err := rows.Scan(
			&testCase.ID,
			&testCase.ProjectID,
			&testCase.SuiteID,
			&testCase.Title,
			&testCase.Description,
			&testCase.Preconditions,
			&testCase.Status,
			&testCase.Priority,
			&testCase.CreatedBy,
			&testCase.UpdatedBy,
			&testCase.Version,
			&testCase.CreatedAt,
			&testCase.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan test case: %v", err)
		}

		// Get steps for each test case
		steps, err := r.GetSteps(testCase.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get test steps: %v", err)
		}
		testCase.Steps = steps

		testCases = append(testCases, testCase)
	}

	return testCases, nil
}

func (r *TestCaseRepository) ListBySuite(suiteID int64) ([]*models.TestCase, error) {
	query := `
		SELECT 
			id, project_id, suite_id, title, description, preconditions,
			status, priority, created_by, updated_by, version,
			created_at, updated_at
		FROM test_cases
		WHERE suite_id = ?
		ORDER BY title`

	rows, err := r.db.Query(query, suiteID)
	if err != nil {
		return nil, fmt.Errorf("failed to list test cases: %v", err)
	}
	defer rows.Close()

	var testCases []*models.TestCase
	for rows.Next() {
		testCase := &models.TestCase{}
		err := rows.Scan(
			&testCase.ID,
			&testCase.ProjectID,
			&testCase.SuiteID,
			&testCase.Title,
			&testCase.Description,
			&testCase.Preconditions,
			&testCase.Status,
			&testCase.Priority,
			&testCase.CreatedBy,
			&testCase.UpdatedBy,
			&testCase.Version,
			&testCase.CreatedAt,
			&testCase.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan test case: %v", err)
		}

		// Get steps for each test case
		steps, err := r.GetSteps(testCase.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get test steps: %v", err)
		}
		testCase.Steps = steps

		testCases = append(testCases, testCase)
	}

	return testCases, nil
}

func (r *TestCaseRepository) GetSteps(testCaseID int64) ([]*models.TestStep, error) {
	query := `
		SELECT 
			id, test_case_id, step_number, step_type, description,
			expected_result, created_at, updated_at
		FROM test_steps
		WHERE test_case_id = ?
		ORDER BY step_number`

	rows, err := r.db.Query(query, testCaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get test steps: %v", err)
	}
	defer rows.Close()

	var steps []*models.TestStep
	for rows.Next() {
		step := &models.TestStep{}
		err := rows.Scan(
			&step.ID,
			&step.TestCaseID,
			&step.StepNumber,
			&step.StepType,
			&step.Description,
			&step.ExpectedResult,
			&step.CreatedAt,
			&step.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan test step: %v", err)
		}

		// Get notes for each step
		notes, err := r.getStepNotes(step.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get step notes: %v", err)
		}
		step.Notes = notes

		// Get attachments for each step
		attachments, err := r.getStepAttachments(step.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get step attachments: %v", err)
		}
		step.Attachments = attachments

		steps = append(steps, step)
	}

	return steps, nil
}

func (r *TestCaseRepository) CreateStep(step *models.TestStep) error {
	now := time.Now()
	query := `
		INSERT INTO test_steps (
			test_case_id, step_number, step_type, description,
			expected_result, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(
		query,
		step.TestCaseID,
		step.StepNumber,
		step.StepType,
		step.Description,
		step.ExpectedResult,
		now,
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to create test step: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %v", err)
	}

	step.ID = id
	step.CreatedAt = now
	step.UpdatedAt = now

	return nil
}

func (r *TestCaseRepository) UpdateStep(step *models.TestStep) error {
	now := time.Now()
	query := `
		UPDATE test_steps SET
			step_number = ?,
			step_type = ?,
			description = ?,
			expected_result = ?,
			updated_at = ?
		WHERE id = ?`

	result, err := r.db.Exec(
		query,
		step.StepNumber,
		step.StepType,
		step.Description,
		step.ExpectedResult,
		now,
		step.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update test step: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("test step not found")
	}

	step.UpdatedAt = now
	return nil
}

func (r *TestCaseRepository) DeleteStep(stepID int64) error {
	result, err := r.db.Exec("DELETE FROM test_steps WHERE id = ?", stepID)
	if err != nil {
		return fmt.Errorf("failed to delete test step: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("test step not found")
	}

	return nil
}

func (r *TestCaseRepository) getStepNotes(stepID int64) ([]*models.StepNote, error) {
	query := `
		SELECT id, step_id, note_text, created_by, created_at
		FROM step_notes
		WHERE step_id = ?
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, stepID)
	if err != nil {
		return nil, fmt.Errorf("failed to get step notes: %v", err)
	}
	defer rows.Close()

	var notes []*models.StepNote
	for rows.Next() {
		note := &models.StepNote{}
		err := rows.Scan(
			&note.ID,
			&note.StepID,
			&note.Content,
			&note.CreatedBy,
			&note.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan step note: %v", err)
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func (r *TestCaseRepository) CreateStepNote(note *models.StepNote) error {
	now := time.Now()
	query := `
		INSERT INTO step_notes (
			step_id, note_text, created_by, created_at
		) VALUES (?, ?, ?, ?)`

	result, err := r.db.Exec(
		query,
		note.StepID,
		note.Content,
		note.CreatedBy,
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to create step note: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %v", err)
	}

	note.ID = id
	note.CreatedAt = now

	return nil
}

func (r *TestCaseRepository) DeleteStepNote(noteID int64) error {
	result, err := r.db.Exec("DELETE FROM step_notes WHERE id = ?", noteID)
	if err != nil {
		return fmt.Errorf("failed to delete step note: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("step note not found")
	}

	return nil
}

func (r *TestCaseRepository) getStepAttachments(stepID int64) ([]*models.StepAttachment, error) {
	query := `
		SELECT id, step_id, file_name, file_path, file_type, file_size, created_by, created_at
		FROM step_attachments
		WHERE step_id = ?
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, stepID)
	if err != nil {
		return nil, fmt.Errorf("failed to get step attachments: %v", err)
	}
	defer rows.Close()

	var attachments []*models.StepAttachment
	for rows.Next() {
		attachment := &models.StepAttachment{}
		err := rows.Scan(
			&attachment.ID,
			&attachment.StepID,
			&attachment.FileName,
			&attachment.FilePath,
			&attachment.FileType,
			&attachment.FileSize,
			&attachment.CreatedBy,
			&attachment.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan step attachment: %v", err)
		}
		attachments = append(attachments, attachment)
	}

	return attachments, nil
}

func (r *TestCaseRepository) CreateStepAttachment(attachment *models.StepAttachment) error {
	now := time.Now()
	query := `
		INSERT INTO step_attachments (
			step_id, file_name, file_path, file_type, file_size, created_by, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(
		query,
		attachment.StepID,
		attachment.FileName,
		attachment.FilePath,
		attachment.FileType,
		attachment.FileSize,
		attachment.CreatedBy,
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to create step attachment: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %v", err)
	}

	attachment.ID = id
	attachment.CreatedAt = now

	return nil
}

func (r *TestCaseRepository) DeleteStepAttachment(attachmentID int64) error {
	result, err := r.db.Exec("DELETE FROM step_attachments WHERE id = ?", attachmentID)
	if err != nil {
		return fmt.Errorf("failed to delete step attachment: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("step attachment not found")
	}

	return nil
}

func (r *TestCaseRepository) GetStepByID(stepID int64) (*models.TestStep, error) {
	query := `
		SELECT 
			id, test_case_id, step_number, step_type, description,
			expected_result, created_at, updated_at
		FROM test_steps
		WHERE id = ?`

	step := &models.TestStep{}
	err := r.db.QueryRow(query, stepID).Scan(
		&step.ID,
		&step.TestCaseID,
		&step.StepNumber,
		&step.StepType,
		&step.Description,
		&step.ExpectedResult,
		&step.CreatedAt,
		&step.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("test step not found")
		}
		return nil, fmt.Errorf("failed to get test step: %v", err)
	}

	// Get notes for the step
	notes, err := r.getStepNotes(step.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get step notes: %v", err)
	}
	step.Notes = notes

	// Get attachments for the step
	attachments, err := r.getStepAttachments(step.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get step attachments: %v", err)
	}
	step.Attachments = attachments

	return step, nil
}

func (r *TestCaseRepository) GetStepAttachmentByID(attachmentID int64) (*models.StepAttachment, error) {
	query := `
		SELECT id, step_id, file_name, file_path, file_type, file_size, created_by, created_at
		FROM step_attachments
		WHERE id = ?`

	attachment := &models.StepAttachment{}
	err := r.db.QueryRow(query, attachmentID).Scan(
		&attachment.ID,
		&attachment.StepID,
		&attachment.FileName,
		&attachment.FilePath,
		&attachment.FileType,
		&attachment.FileSize,
		&attachment.CreatedBy,
		&attachment.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("step attachment not found")
		}
		return nil, fmt.Errorf("failed to get step attachment: %v", err)
	}

	return attachment, nil
}
