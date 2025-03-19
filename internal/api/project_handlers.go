package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mihaamiharu/test-case-management-be/internal/models"
	"github.com/mihaamiharu/test-case-management-be/internal/repository"
	"github.com/mihaamiharu/test-case-management-be/internal/services"
)

// ProjectHandler handles project-related API endpoints
type ProjectHandler struct {
	projectService       *services.ProjectService
	projectAccessService *services.ProjectAccessService
}

// NewProjectHandler creates a new project handler
func NewProjectHandler(projectService *services.ProjectService, projectAccessService *services.ProjectAccessService) *ProjectHandler {
	return &ProjectHandler{
		projectService:       projectService,
		projectAccessService: projectAccessService,
	}
}

// CreateProject handles creating a new project
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var projectCreate models.ProjectCreate
	if err := c.ShouldBindJSON(&projectCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	userModel := user.(*models.User)

	project, err := h.projectService.Create(&projectCreate, userModel.ID)
	if err != nil {
		if err == repository.ErrProjectExists {
			c.JSON(http.StatusConflict, gin.H{"error": "Project with this name already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}

	c.JSON(http.StatusCreated, project.ToResponse())
}

// GetProject handles retrieving a project by ID
func (h *ProjectHandler) GetProject(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	userModel := user.(*models.User)

	// Check if user has access to the project
	hasAccess, err := h.projectAccessService.HasViewAccess(id, userModel.ID)
	if err != nil {
		if err == repository.ErrProjectNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check project access"})
		return
	}

	// Only allow users with access to view the project
	if !hasAccess && userModel.Role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to view this project"})
		return
	}

	project, err := h.projectService.GetByID(id)
	if err != nil {
		if err == repository.ErrProjectNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve project"})
		return
	}

	c.JSON(http.StatusOK, project.ToResponse())
}

// UpdateProject handles updating a project
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var projectUpdate models.ProjectUpdate
	if err := c.ShouldBindJSON(&projectUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	userModel := user.(*models.User)

	// Check if user has edit access to the project
	hasEditAccess, err := h.projectAccessService.HasEditAccess(id, userModel.ID)
	if err != nil {
		if err == repository.ErrProjectNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check project access"})
		return
	}

	// Only allow users with edit access to update the project (unless admin)
	if !hasEditAccess && userModel.Role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this project"})
		return
	}

	project, err := h.projectService.Update(id, &projectUpdate)
	if err != nil {
		if err == repository.ErrProjectNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project"})
		return
	}

	c.JSON(http.StatusOK, project.ToResponse())
}

// DeleteProject handles deleting a project
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	userModel := user.(*models.User)

	// Check if user is the owner of the project
	isOwner, err := h.projectService.IsOwner(id, userModel.ID)
	if err != nil {
		if err == repository.ErrProjectNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check project ownership"})
		return
	}

	// Only allow the owner to delete the project (unless admin)
	if !isOwner && userModel.Role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this project"})
		return
	}

	if err := h.projectService.Delete(id); err != nil {
		if err == repository.ErrProjectNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}

// ListProjects handles retrieving all projects with pagination
func (h *ProjectHandler) ListProjects(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	userModel := user.(*models.User)

	var projects []*models.Project
	var err error

	// If admin, show all projects, otherwise show only projects the user has access to
	if userModel.Role == models.RoleAdmin {
		projects, err = h.projectService.List(page, pageSize)
	} else {
		projects, err = h.projectAccessService.GetAccessibleProjects(userModel.ID, page, pageSize)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve projects"})
		return
	}

	// Convert to response objects
	var responses []models.ProjectResponse
	for _, project := range projects {
		responses = append(responses, project.ToResponse())
	}

	c.JSON(http.StatusOK, responses)
}
