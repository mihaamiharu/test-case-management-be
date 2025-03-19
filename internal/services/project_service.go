package services

import (
	"github.com/mihaamiharu/test-case-management-be/internal/models"
	"github.com/mihaamiharu/test-case-management-be/internal/repository"
)

// ProjectService handles business logic for projects
type ProjectService struct {
	projectRepo repository.ProjectRepositoryInterface
}

// NewProjectService creates a new project service
func NewProjectService(projectRepo repository.ProjectRepositoryInterface) *ProjectService {
	return &ProjectService{
		projectRepo: projectRepo,
	}
}

// Create creates a new project
func (s *ProjectService) Create(projectCreate *models.ProjectCreate, ownerID int64) (*models.Project, error) {
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

// GetByID retrieves a project by ID
func (s *ProjectService) GetByID(id int64) (*models.Project, error) {
	return s.projectRepo.GetByID(id)
}

// Update updates an existing project
func (s *ProjectService) Update(id int64, projectUpdate *models.ProjectUpdate) (*models.Project, error) {
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

// Delete removes a project
func (s *ProjectService) Delete(id int64) error {
	return s.projectRepo.Delete(id)
}

// ListByOwner retrieves all projects for a specific owner
func (s *ProjectService) ListByOwner(ownerID int64) ([]*models.Project, error) {
	return s.projectRepo.ListByOwner(ownerID)
}

// List retrieves all projects with pagination
func (s *ProjectService) List(page, pageSize int) ([]*models.Project, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	return s.projectRepo.List(pageSize, offset)
}

// IsOwner checks if a user is the owner of a project
func (s *ProjectService) IsOwner(projectID, userID int64) (bool, error) {
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return false, err
	}

	return project.OwnerID == userID, nil
}
