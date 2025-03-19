package services

import (
	"errors"
	"testing"
	"time"

	"github.com/mihaamiharu/test-case-management-be/internal/models"
	"github.com/mihaamiharu/test-case-management-be/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ProjectRepositoryInterface defines the interface for project repository
type ProjectRepositoryInterface interface {
	Create(project *models.Project) error
	GetByID(id int64) (*models.Project, error)
	Update(project *models.Project) error
	Delete(id int64) error
	ListByOwner(ownerID int64) ([]*models.Project, error)
	List(limit, offset int) ([]*models.Project, error)
}

// MockProjectRepository is a mock implementation of the project repository
type MockProjectRepository struct {
	mock.Mock
}

func (m *MockProjectRepository) Create(project *models.Project) error {
	args := m.Called(project)
	return args.Error(0)
}

func (m *MockProjectRepository) GetByID(id int64) (*models.Project, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Project), args.Error(1)
}

func (m *MockProjectRepository) Update(project *models.Project) error {
	args := m.Called(project)
	return args.Error(0)
}

func (m *MockProjectRepository) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockProjectRepository) ListByOwner(ownerID int64) ([]*models.Project, error) {
	args := m.Called(ownerID)
	return args.Get(0).([]*models.Project), args.Error(1)
}

func (m *MockProjectRepository) List(limit, offset int) ([]*models.Project, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]*models.Project), args.Error(1)
}

// ProjectService for testing
type testProjectService struct {
	projectRepo *MockProjectRepository
}

func (s *testProjectService) Create(projectCreate *models.ProjectCreate, ownerID int64) (*models.Project, error) {
	project := &models.Project{
		Name:        projectCreate.Name,
		Description: projectCreate.Description,
		OwnerID:     ownerID,
	}

	if err := s.projectRepo.Create(project); err != nil {
		return nil, err
	}

	return project, nil
}

func (s *testProjectService) GetByID(id int64) (*models.Project, error) {
	return s.projectRepo.GetByID(id)
}

func (s *testProjectService) Update(id int64, projectUpdate *models.ProjectUpdate) (*models.Project, error) {
	project, err := s.projectRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if projectUpdate.Name != "" {
		project.Name = projectUpdate.Name
	}

	// Description can be empty, so we always update it
	project.Description = projectUpdate.Description

	if err := s.projectRepo.Update(project); err != nil {
		return nil, err
	}

	return project, nil
}

func (s *testProjectService) Delete(id int64) error {
	return s.projectRepo.Delete(id)
}

func (s *testProjectService) ListByOwner(ownerID int64) ([]*models.Project, error) {
	return s.projectRepo.ListByOwner(ownerID)
}

func (s *testProjectService) List(page, pageSize int) ([]*models.Project, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	return s.projectRepo.List(pageSize, offset)
}

func (s *testProjectService) IsOwner(projectID, userID int64) (bool, error) {
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return false, err
	}

	return project.OwnerID == userID, nil
}

// NewTestProjectService creates a new test project service
func NewTestProjectService(mockRepo *MockProjectRepository) *testProjectService {
	return &testProjectService{
		projectRepo: mockRepo,
	}
}

func TestProjectService_Create(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	service := NewTestProjectService(mockRepo)
	now := time.Now()

	// Test case: successful creation
	t.Run("Success", func(t *testing.T) {
		projectCreate := &models.ProjectCreate{
			Name:        "Test Project",
			Description: "Test Description",
		}

		// Setup expectations
		mockRepo.On("Create", mock.MatchedBy(func(p *models.Project) bool {
			return p.Name == projectCreate.Name && p.Description == projectCreate.Description && p.OwnerID == int64(1)
		})).Run(func(args mock.Arguments) {
			project := args.Get(0).(*models.Project)
			project.ID = 1
			project.CreatedAt = now
			project.UpdatedAt = now
		}).Return(nil).Once()

		// Execute
		project, err := service.Create(projectCreate, 1)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, project)
		assert.Equal(t, int64(1), project.ID)
		assert.Equal(t, "Test Project", project.Name)
		assert.Equal(t, "Test Description", project.Description)
		assert.Equal(t, int64(1), project.OwnerID)
		assert.Equal(t, now, project.CreatedAt)
		assert.Equal(t, now, project.UpdatedAt)
		mockRepo.AssertExpectations(t)
	})

	// Test case: repository error
	t.Run("RepositoryError", func(t *testing.T) {
		projectCreate := &models.ProjectCreate{
			Name:        "Test Project",
			Description: "Test Description",
		}

		// Setup expectations
		mockRepo.On("Create", mock.MatchedBy(func(p *models.Project) bool {
			return p.Name == projectCreate.Name && p.Description == projectCreate.Description && p.OwnerID == int64(1)
		})).Return(repository.ErrProjectExists).Once()

		// Execute
		project, err := service.Create(projectCreate, 1)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, repository.ErrProjectExists, err)
		assert.Nil(t, project)
		mockRepo.AssertExpectations(t)
	})
}

func TestProjectService_GetByID(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	service := NewTestProjectService(mockRepo)
	now := time.Now()

	// Test case: project found
	t.Run("Success", func(t *testing.T) {
		expectedProject := &models.Project{
			ID:          1,
			Name:        "Test Project",
			Description: "Test Description",
			OwnerID:     1,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		// Setup expectations
		mockRepo.On("GetByID", int64(1)).Return(expectedProject, nil).Once()

		// Execute
		project, err := service.GetByID(1)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedProject, project)
		mockRepo.AssertExpectations(t)
	})

	// Test case: project not found
	t.Run("NotFound", func(t *testing.T) {
		// Setup expectations
		mockRepo.On("GetByID", int64(2)).Return(nil, repository.ErrProjectNotFound).Once()

		// Execute
		project, err := service.GetByID(2)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, repository.ErrProjectNotFound, err)
		assert.Nil(t, project)
		mockRepo.AssertExpectations(t)
	})
}

func TestProjectService_Update(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	service := NewTestProjectService(mockRepo)
	now := time.Now()

	// Test case: successful update
	t.Run("Success", func(t *testing.T) {
		existingProject := &models.Project{
			ID:          1,
			Name:        "Old Name",
			Description: "Old Description",
			OwnerID:     1,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		projectUpdate := &models.ProjectUpdate{
			Name:        "New Name",
			Description: "New Description",
		}

		updatedProject := &models.Project{
			ID:          1,
			Name:        "New Name",
			Description: "New Description",
			OwnerID:     1,
			CreatedAt:   now,
			UpdatedAt:   now.Add(time.Hour),
		}

		// Setup expectations
		mockRepo.On("GetByID", int64(1)).Return(existingProject, nil).Once()
		mockRepo.On("Update", mock.MatchedBy(func(p *models.Project) bool {
			return p.ID == existingProject.ID && p.Name == projectUpdate.Name && p.Description == projectUpdate.Description
		})).Run(func(args mock.Arguments) {
			project := args.Get(0).(*models.Project)
			project.UpdatedAt = now.Add(time.Hour)
		}).Return(nil).Once()

		// Execute
		project, err := service.Update(1, projectUpdate)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, updatedProject.Name, project.Name)
		assert.Equal(t, updatedProject.Description, project.Description)
		assert.Equal(t, updatedProject.UpdatedAt, project.UpdatedAt)
		mockRepo.AssertExpectations(t)
	})

	// Test case: project not found
	t.Run("NotFound", func(t *testing.T) {
		projectUpdate := &models.ProjectUpdate{
			Name:        "New Name",
			Description: "New Description",
		}

		// Setup expectations
		mockRepo.On("GetByID", int64(2)).Return(nil, repository.ErrProjectNotFound).Once()

		// Execute
		project, err := service.Update(2, projectUpdate)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, repository.ErrProjectNotFound, err)
		assert.Nil(t, project)
		mockRepo.AssertExpectations(t)
	})

	// Test case: update error
	t.Run("UpdateError", func(t *testing.T) {
		existingProject := &models.Project{
			ID:          1,
			Name:        "Old Name",
			Description: "Old Description",
			OwnerID:     1,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		projectUpdate := &models.ProjectUpdate{
			Name:        "New Name",
			Description: "New Description",
		}

		// Setup expectations
		mockRepo.On("GetByID", int64(1)).Return(existingProject, nil).Once()
		mockRepo.On("Update", mock.MatchedBy(func(p *models.Project) bool {
			return p.ID == existingProject.ID && p.Name == projectUpdate.Name && p.Description == projectUpdate.Description
		})).Return(errors.New("update error")).Once()

		// Execute
		project, err := service.Update(1, projectUpdate)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, project)
		mockRepo.AssertExpectations(t)
	})
}

func TestProjectService_Delete(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	service := NewTestProjectService(mockRepo)

	// Test case: successful delete
	t.Run("Success", func(t *testing.T) {
		// Setup expectations
		mockRepo.On("Delete", int64(1)).Return(nil).Once()

		// Execute
		err := service.Delete(1)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	// Test case: project not found
	t.Run("NotFound", func(t *testing.T) {
		// Setup expectations
		mockRepo.On("Delete", int64(2)).Return(repository.ErrProjectNotFound).Once()

		// Execute
		err := service.Delete(2)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, repository.ErrProjectNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestProjectService_ListByOwner(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	service := NewTestProjectService(mockRepo)
	now := time.Now()

	// Test case: successful list
	t.Run("Success", func(t *testing.T) {
		expectedProjects := []*models.Project{
			{
				ID:          1,
				Name:        "Project 1",
				Description: "Description 1",
				OwnerID:     1,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			{
				ID:          2,
				Name:        "Project 2",
				Description: "Description 2",
				OwnerID:     1,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		}

		// Setup expectations
		mockRepo.On("ListByOwner", int64(1)).Return(expectedProjects, nil).Once()

		// Execute
		projects, err := service.ListByOwner(1)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedProjects, projects)
		mockRepo.AssertExpectations(t)
	})

	// Test case: no projects
	t.Run("NoProjects", func(t *testing.T) {
		// Setup expectations
		mockRepo.On("ListByOwner", int64(2)).Return([]*models.Project{}, nil).Once()

		// Execute
		projects, err := service.ListByOwner(2)

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, projects)
		mockRepo.AssertExpectations(t)
	})

	// Test case: repository error
	t.Run("RepositoryError", func(t *testing.T) {
		// Setup expectations
		mockRepo.On("ListByOwner", int64(3)).Return([]*models.Project{}, errors.New("repository error")).Once()

		// Execute
		projects, err := service.ListByOwner(3)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, projects)
		mockRepo.AssertExpectations(t)
	})
}

func TestProjectService_List(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	service := NewTestProjectService(mockRepo)
	now := time.Now()

	// Test case: successful list
	t.Run("Success", func(t *testing.T) {
		expectedProjects := []*models.Project{
			{
				ID:          1,
				Name:        "Project 1",
				Description: "Description 1",
				OwnerID:     1,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			{
				ID:          2,
				Name:        "Project 2",
				Description: "Description 2",
				OwnerID:     2,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		}

		// Setup expectations
		mockRepo.On("List", 10, 0).Return(expectedProjects, nil).Once()

		// Execute
		projects, err := service.List(1, 10)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedProjects, projects)
		mockRepo.AssertExpectations(t)
	})

	// Test case: no projects
	t.Run("NoProjects", func(t *testing.T) {
		// Setup expectations
		mockRepo.On("List", 10, 10).Return([]*models.Project{}, nil).Once()

		// Execute
		projects, err := service.List(2, 10)

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, projects)
		mockRepo.AssertExpectations(t)
	})

	// Test case: repository error
	t.Run("RepositoryError", func(t *testing.T) {
		// Setup expectations
		mockRepo.On("List", 10, 20).Return([]*models.Project{}, errors.New("repository error")).Once()

		// Execute
		projects, err := service.List(3, 10)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, projects)
		mockRepo.AssertExpectations(t)
	})

	// Test case: negative page
	t.Run("NegativePage", func(t *testing.T) {
		expectedProjects := []*models.Project{
			{
				ID:          1,
				Name:        "Project 1",
				Description: "Description 1",
				OwnerID:     1,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		}

		// Setup expectations
		mockRepo.On("List", 10, 0).Return(expectedProjects, nil).Once()

		// Execute
		projects, err := service.List(-1, 10)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedProjects, projects)
		mockRepo.AssertExpectations(t)
	})

	// Test case: negative page size
	t.Run("NegativePageSize", func(t *testing.T) {
		expectedProjects := []*models.Project{
			{
				ID:          1,
				Name:        "Project 1",
				Description: "Description 1",
				OwnerID:     1,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		}

		// Setup expectations
		mockRepo.On("List", 10, 0).Return(expectedProjects, nil).Once()

		// Execute
		projects, err := service.List(1, -5)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedProjects, projects)
		mockRepo.AssertExpectations(t)
	})
}

func TestProjectService_IsOwner(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	service := NewTestProjectService(mockRepo)
	now := time.Now()

	// Test case: user is owner
	t.Run("IsOwner", func(t *testing.T) {
		project := &models.Project{
			ID:          1,
			Name:        "Test Project",
			Description: "Test Description",
			OwnerID:     1,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		// Setup expectations
		mockRepo.On("GetByID", int64(1)).Return(project, nil).Once()

		// Execute
		isOwner, err := service.IsOwner(1, 1)

		// Assert
		assert.NoError(t, err)
		assert.True(t, isOwner)
		mockRepo.AssertExpectations(t)
	})

	// Test case: user is not owner
	t.Run("NotOwner", func(t *testing.T) {
		project := &models.Project{
			ID:          1,
			Name:        "Test Project",
			Description: "Test Description",
			OwnerID:     1,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		// Setup expectations
		mockRepo.On("GetByID", int64(1)).Return(project, nil).Once()

		// Execute
		isOwner, err := service.IsOwner(1, 2)

		// Assert
		assert.NoError(t, err)
		assert.False(t, isOwner)
		mockRepo.AssertExpectations(t)
	})

	// Test case: project not found
	t.Run("ProjectNotFound", func(t *testing.T) {
		// Setup expectations
		mockRepo.On("GetByID", int64(2)).Return(nil, repository.ErrProjectNotFound).Once()

		// Execute
		isOwner, err := service.IsOwner(2, 1)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, repository.ErrProjectNotFound, err)
		assert.False(t, isOwner)
		mockRepo.AssertExpectations(t)
	})
}
