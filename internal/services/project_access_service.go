package services

import (
	"github.com/mihaamiharu/test-case-management-be/internal/models"
	"github.com/mihaamiharu/test-case-management-be/internal/repository"
)

// ProjectAccessService handles business logic for project access
type ProjectAccessService struct {
	projectAccessRepo repository.ProjectAccessRepositoryInterface
	projectRepo       repository.ProjectRepositoryInterface
	userRepo          repository.UserRepositoryInterface
}

// NewProjectAccessService creates a new project access service
func NewProjectAccessService(
	projectAccessRepo repository.ProjectAccessRepositoryInterface,
	projectRepo repository.ProjectRepositoryInterface,
	userRepo repository.UserRepositoryInterface,
) *ProjectAccessService {
	return &ProjectAccessService{
		projectAccessRepo: projectAccessRepo,
		projectRepo:       projectRepo,
		userRepo:          userRepo,
	}
}

// GrantAccess grants access to a project for a user
func (s *ProjectAccessService) GrantAccess(projectID int64, accessCreate *models.ProjectAccessCreate) (*models.ProjectAccess, error) {
	// Check if project exists
	_, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return nil, err
	}

	// Check if user exists
	_, err = s.userRepo.GetByID(accessCreate.UserID)
	if err != nil {
		return nil, err
	}

	// Create access
	access := &models.ProjectAccess{
		ProjectID: projectID,
		UserID:    accessCreate.UserID,
		Level:     accessCreate.Level,
	}

	if err := s.projectAccessRepo.Create(access); err != nil {
		return nil, err
	}

	return access, nil
}

// UpdateAccess updates a user's access level to a project
func (s *ProjectAccessService) UpdateAccess(accessID int64, accessUpdate *models.ProjectAccessUpdate) (*models.ProjectAccess, error) {
	// Get existing access
	access, err := s.projectAccessRepo.GetByID(accessID)
	if err != nil {
		return nil, err
	}

	// Update access level
	access.Level = accessUpdate.Level

	if err := s.projectAccessRepo.Update(access); err != nil {
		return nil, err
	}

	return access, nil
}

// RevokeAccess removes a user's access to a project
func (s *ProjectAccessService) RevokeAccess(accessID int64) error {
	return s.projectAccessRepo.Delete(accessID)
}

// GetProjectAccess retrieves all access records for a project
func (s *ProjectAccessService) GetProjectAccess(projectID int64) ([]*models.ProjectAccess, error) {
	return s.projectAccessRepo.ListByProject(projectID)
}

// GetUserAccess retrieves all access records for a user
func (s *ProjectAccessService) GetUserAccess(userID int64) ([]*models.ProjectAccess, error) {
	return s.projectAccessRepo.ListByUser(userID)
}

// HasEditAccess checks if a user has edit access to a project
func (s *ProjectAccessService) HasEditAccess(projectID, userID int64) (bool, error) {
	// Project owner always has edit access
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return false, err
	}

	if project.OwnerID == userID {
		return true, nil
	}

	// Check if user has explicit edit access
	return s.projectAccessRepo.HasEditAccess(projectID, userID)
}

// HasViewAccess checks if a user has at least view access to a project
func (s *ProjectAccessService) HasViewAccess(projectID, userID int64) (bool, error) {
	// Project owner always has view access
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return false, err
	}

	if project.OwnerID == userID {
		return true, nil
	}

	// Check if user has explicit view access
	return s.projectAccessRepo.HasViewAccess(projectID, userID)
}

// GetAccessibleProjects retrieves all projects a user has access to
func (s *ProjectAccessService) GetAccessibleProjects(userID int64, page, pageSize int) ([]*models.Project, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	// Get projects owned by the user
	ownedProjects, err := s.projectRepo.ListByOwner(userID)
	if err != nil {
		return nil, err
	}

	// Get project IDs the user has access to
	accessibleProjectIDs, err := s.projectAccessRepo.GetProjectIDsByUserAccess(userID)
	if err != nil {
		return nil, err
	}

	// If user doesn't have access to any projects and doesn't own any, return empty list
	if len(accessibleProjectIDs) == 0 && len(ownedProjects) == 0 {
		return []*models.Project{}, nil
	}

	// Create a map of owned project IDs for quick lookup
	ownedProjectIDs := make(map[int64]bool)
	for _, project := range ownedProjects {
		ownedProjectIDs[project.ID] = true
	}

	// Filter out project IDs that are already owned by the user
	var uniqueAccessibleProjectIDs []int64
	for _, id := range accessibleProjectIDs {
		if !ownedProjectIDs[id] {
			uniqueAccessibleProjectIDs = append(uniqueAccessibleProjectIDs, id)
		}
	}

	// If user doesn't have access to any additional projects, return only owned projects
	if len(uniqueAccessibleProjectIDs) == 0 {
		// Apply pagination to owned projects
		start := (page - 1) * pageSize
		end := start + pageSize
		if start >= len(ownedProjects) {
			return []*models.Project{}, nil
		}
		if end > len(ownedProjects) {
			end = len(ownedProjects)
		}
		return ownedProjects[start:end], nil
	}

	// Get projects the user has access to
	var accessibleProjects []*models.Project
	for _, projectID := range uniqueAccessibleProjectIDs {
		project, err := s.projectRepo.GetByID(projectID)
		if err != nil {
			// Skip projects that can't be retrieved
			continue
		}
		accessibleProjects = append(accessibleProjects, project)
	}

	// Combine owned and accessible projects
	allProjects := append(ownedProjects, accessibleProjects...)

	// Apply pagination
	start := (page - 1) * pageSize
	end := start + pageSize
	if start >= len(allProjects) {
		return []*models.Project{}, nil
	}
	if end > len(allProjects) {
		end = len(allProjects)
	}

	return allProjects[start:end], nil
}
