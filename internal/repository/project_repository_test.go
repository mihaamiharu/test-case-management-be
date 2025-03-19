package repository

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mihaamiharu/test-case-management-be/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestProjectRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewProjectRepository(db)

	// Test case: successful project creation
	t.Run("Success", func(t *testing.T) {
		// Setup expectations
		mock.ExpectQuery("SELECT COUNT").
			WithArgs("Test Project", int64(1)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		mock.ExpectExec("INSERT INTO projects").
			WithArgs("Test Project", "Test Description", int64(1), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Create project
		project := &models.Project{
			Name:        "Test Project",
			Description: "Test Description",
			OwnerID:     1,
		}

		// Execute
		err := repo.Create(project)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int64(1), project.ID)
		assert.NotZero(t, project.CreatedAt)
		assert.NotZero(t, project.UpdatedAt)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Test case: project already exists
	t.Run("ProjectExists", func(t *testing.T) {
		// Setup expectations
		mock.ExpectQuery("SELECT COUNT").
			WithArgs("Existing Project", int64(1)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		// Create project
		project := &models.Project{
			Name:        "Existing Project",
			Description: "Test Description",
			OwnerID:     1,
		}

		// Execute
		err := repo.Create(project)

		// Assert
		assert.Equal(t, ErrProjectExists, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Test case: database error
	t.Run("DatabaseError", func(t *testing.T) {
		// Setup expectations
		mock.ExpectQuery("SELECT COUNT").
			WithArgs("Test Project", int64(1)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		mock.ExpectExec("INSERT INTO projects").
			WithArgs("Test Project", "Test Description", int64(1), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(sql.ErrConnDone)

		// Create project
		project := &models.Project{
			Name:        "Test Project",
			Description: "Test Description",
			OwnerID:     1,
		}

		// Execute
		err := repo.Create(project)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, sql.ErrConnDone, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestProjectRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewProjectRepository(db)
	now := time.Now()

	// Test case: project found
	t.Run("Success", func(t *testing.T) {
		// Setup expectations
		rows := sqlmock.NewRows([]string{"id", "name", "description", "owner_id", "created_at", "updated_at"}).
			AddRow(1, "Test Project", "Test Description", 1, now, now)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, description, owner_id, created_at, updated_at FROM projects WHERE id = ?")).
			WithArgs(1).
			WillReturnRows(rows)

		// Execute
		project, err := repo.GetByID(1)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, project)
		assert.Equal(t, int64(1), project.ID)
		assert.Equal(t, "Test Project", project.Name)
		assert.Equal(t, "Test Description", project.Description)
		assert.Equal(t, int64(1), project.OwnerID)
		assert.Equal(t, now, project.CreatedAt)
		assert.Equal(t, now, project.UpdatedAt)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Test case: project not found
	t.Run("NotFound", func(t *testing.T) {
		// Setup expectations
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, description, owner_id, created_at, updated_at FROM projects WHERE id = ?")).
			WithArgs(2).
			WillReturnError(sql.ErrNoRows)

		// Execute
		project, err := repo.GetByID(2)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, ErrProjectNotFound, err)
		assert.Nil(t, project)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestProjectRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewProjectRepository(db)
	now := time.Now()

	// Test case: successful update
	t.Run("Success", func(t *testing.T) {
		// Setup expectations for GetByID
		rows := sqlmock.NewRows([]string{"id", "name", "description", "owner_id", "created_at", "updated_at"}).
			AddRow(1, "Old Name", "Old Description", 1, now, now)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, description, owner_id, created_at, updated_at FROM projects WHERE id = ?")).
			WithArgs(1).
			WillReturnRows(rows)

		// Setup expectations for Update
		mock.ExpectExec(regexp.QuoteMeta("UPDATE projects SET name = ?, description = ?, updated_at = ? WHERE id = ?")).
			WithArgs("New Name", "New Description", sqlmock.AnyArg(), 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Create project
		project := &models.Project{
			ID:          1,
			Name:        "New Name",
			Description: "New Description",
			OwnerID:     1,
		}

		// Execute
		err := repo.Update(project)

		// Assert
		assert.NoError(t, err)
		assert.NotEqual(t, now, project.UpdatedAt)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Test case: project not found
	t.Run("NotFound", func(t *testing.T) {
		// Setup expectations
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, description, owner_id, created_at, updated_at FROM projects WHERE id = ?")).
			WithArgs(2).
			WillReturnError(sql.ErrNoRows)

		// Create project
		project := &models.Project{
			ID:          2,
			Name:        "New Name",
			Description: "New Description",
			OwnerID:     1,
		}

		// Execute
		err := repo.Update(project)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, ErrProjectNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestProjectRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewProjectRepository(db)
	now := time.Now()

	// Test case: successful delete
	t.Run("Success", func(t *testing.T) {
		// Setup expectations for GetByID
		rows := sqlmock.NewRows([]string{"id", "name", "description", "owner_id", "created_at", "updated_at"}).
			AddRow(1, "Test Project", "Test Description", 1, now, now)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, description, owner_id, created_at, updated_at FROM projects WHERE id = ?")).
			WithArgs(1).
			WillReturnRows(rows)

		// Setup expectations for Delete
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM projects WHERE id = ?")).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Execute
		err := repo.Delete(1)

		// Assert
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Test case: project not found
	t.Run("NotFound", func(t *testing.T) {
		// Setup expectations
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, description, owner_id, created_at, updated_at FROM projects WHERE id = ?")).
			WithArgs(2).
			WillReturnError(sql.ErrNoRows)

		// Execute
		err := repo.Delete(2)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, ErrProjectNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestProjectRepository_ListByOwner(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewProjectRepository(db)
	now := time.Now()

	// Test case: successful list
	t.Run("Success", func(t *testing.T) {
		// Setup expectations
		rows := sqlmock.NewRows([]string{"id", "name", "description", "owner_id", "created_at", "updated_at"}).
			AddRow(1, "Project 1", "Description 1", 1, now, now).
			AddRow(2, "Project 2", "Description 2", 1, now, now)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, description, owner_id, created_at, updated_at FROM projects WHERE owner_id = ? ORDER BY created_at DESC")).
			WithArgs(1).
			WillReturnRows(rows)

		// Execute
		projects, err := repo.ListByOwner(1)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, projects, 2)
		assert.Equal(t, int64(1), projects[0].ID)
		assert.Equal(t, "Project 1", projects[0].Name)
		assert.Equal(t, int64(2), projects[1].ID)
		assert.Equal(t, "Project 2", projects[1].Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Test case: no projects found
	t.Run("NoProjects", func(t *testing.T) {
		// Setup expectations
		rows := sqlmock.NewRows([]string{"id", "name", "description", "owner_id", "created_at", "updated_at"})

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, description, owner_id, created_at, updated_at FROM projects WHERE owner_id = ? ORDER BY created_at DESC")).
			WithArgs(2).
			WillReturnRows(rows)

		// Execute
		projects, err := repo.ListByOwner(2)

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, projects)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
