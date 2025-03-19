package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mihaamiharu/test-case-management-be/internal/models"
	"github.com/mihaamiharu/test-case-management-be/internal/repository"
	"github.com/mihaamiharu/test-case-management-be/internal/services"
)

// ProjectAccessHandler handles project access-related API endpoints
type ProjectAccessHandler struct {
	projectAccessService *services.ProjectAccessService
	projectService       *services.ProjectService
}

// NewProjectAccessHandler creates a new project access handler
func NewProjectAccessHandler(projectAccessService *services.ProjectAccessService, projectService *services.ProjectService) *ProjectAccessHandler {
	return &ProjectAccessHandler{
		projectAccessService: projectAccessService,
		projectService:       projectService,
	}
}

// GrantAccess handles granting access to a project
func (h *ProjectAccessHandler) GrantAccess(c *gin.Context) {
	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var accessCreate models.ProjectAccessCreate
	if err := c.ShouldBindJSON(&accessCreate); err != nil {
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

	// Check if user is the owner of the project
	isOwner, err := h.projectService.IsOwner(projectID, userModel.ID)
	if err != nil {
		if err == repository.ErrProjectNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check project ownership"})
		return
	}

	// Only allow the owner to grant access to the project
	if !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to grant access to this project"})
		return
	}

	// Don't allow granting access to the owner
	if accessCreate.UserID == userModel.ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You can't grant access to yourself as the owner"})
		return
	}

	access, err := h.projectAccessService.GrantAccess(projectID, &accessCreate)
	if err != nil {
		if err == repository.ErrProjectAccessExists {
			c.JSON(http.StatusConflict, gin.H{"error": "User already has access to this project"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to grant access"})
		return
	}

	c.JSON(http.StatusCreated, access.ToResponse())
}

// UpdateAccess handles updating access to a project
func (h *ProjectAccessHandler) UpdateAccess(c *gin.Context) {
	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	accessID, err := strconv.ParseInt(c.Param("accessId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid access ID"})
		return
	}

	var accessUpdate models.ProjectAccessUpdate
	if err := c.ShouldBindJSON(&accessUpdate); err != nil {
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

	// Check if user is the owner of the project
	isOwner, err := h.projectService.IsOwner(projectID, userModel.ID)
	if err != nil {
		if err == repository.ErrProjectNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check project ownership"})
		return
	}

	// Only allow the owner to update access to the project
	if !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update access to this project"})
		return
	}

	access, err := h.projectAccessService.UpdateAccess(accessID, &accessUpdate)
	if err != nil {
		if err == repository.ErrProjectAccessNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Access record not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update access"})
		return
	}

	c.JSON(http.StatusOK, access.ToResponse())
}

// RevokeAccess handles revoking access to a project
func (h *ProjectAccessHandler) RevokeAccess(c *gin.Context) {
	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	accessID, err := strconv.ParseInt(c.Param("accessId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid access ID"})
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
	isOwner, err := h.projectService.IsOwner(projectID, userModel.ID)
	if err != nil {
		if err == repository.ErrProjectNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check project ownership"})
		return
	}

	// Only allow the owner to revoke access to the project
	if !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to revoke access to this project"})
		return
	}

	if err := h.projectAccessService.RevokeAccess(accessID); err != nil {
		if err == repository.ErrProjectAccessNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Access record not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke access"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Access revoked successfully"})
}

// ListAccess handles listing all access records for a project
func (h *ProjectAccessHandler) ListAccess(c *gin.Context) {
	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
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
	isOwner, err := h.projectService.IsOwner(projectID, userModel.ID)
	if err != nil {
		if err == repository.ErrProjectNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check project ownership"})
		return
	}

	// Only allow the owner to list access records for the project
	if !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to view access records for this project"})
		return
	}

	accessList, err := h.projectAccessService.GetProjectAccess(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve access records"})
		return
	}

	// Convert to response objects
	var responses []models.ProjectAccessResponse
	for _, access := range accessList {
		responses = append(responses, access.ToResponse())
	}

	c.JSON(http.StatusOK, responses)
}
