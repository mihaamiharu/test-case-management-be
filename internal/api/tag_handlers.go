package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mihaamiharu/test-case-management-be/internal/models"
	"github.com/mihaamiharu/test-case-management-be/internal/repository"
	"github.com/mihaamiharu/test-case-management-be/internal/service"
)

// TagHandler handles tag related requests
type TagHandler struct {
	tagService *service.TagService
}

// NewTagHandler creates a new tag handler
func NewTagHandler(tagService *service.TagService) *TagHandler {
	return &TagHandler{
		tagService: tagService,
	}
}

// ListTags handles listing all tags
func (h *TagHandler) ListTags(c *gin.Context) {
	tags, err := h.tagService.ListTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tags"})
		return
	}

	// Convert to response objects
	responses := make([]models.TagResponse, 0, len(tags))
	for _, tag := range tags {
		responses = append(responses, *tag.ToResponse())
	}

	c.JSON(http.StatusOK, responses)
}

// CreateTag handles creating a new tag
func (h *TagHandler) CreateTag(c *gin.Context) {
	var tagCreate models.TagCreate
	if err := c.ShouldBindJSON(&tagCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create tag from request data
	tag := &models.Tag{
		Name: tagCreate.Name,
	}

	err := h.tagService.CreateTag(tag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tag"})
		return
	}

	c.JSON(http.StatusCreated, tag.ToResponse())
}

// GetTag handles retrieving a tag by ID
func (h *TagHandler) GetTag(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	tag, err := h.tagService.GetTagByID(id)
	if err != nil {
		if err == repository.ErrTagNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tag"})
		return
	}

	c.JSON(http.StatusOK, tag.ToResponse())
}

// DeleteTag handles deleting a tag
func (h *TagHandler) DeleteTag(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	err = h.tagService.DeleteTag(id)
	if err != nil {
		if err == repository.ErrTagNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete tag"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag deleted successfully"})
}
