package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mihaamiharu/test-case-management-be/internal/models"
	"github.com/mihaamiharu/test-case-management-be/internal/repository"
	"github.com/mihaamiharu/test-case-management-be/internal/service"
)

// TestSuiteHandler handles test suite related requests
type TestSuiteHandler struct {
	testSuiteService *service.TestSuiteService
}

// NewTestSuiteHandler creates a new test suite handler
func NewTestSuiteHandler(testSuiteService *service.TestSuiteService) *TestSuiteHandler {
	return &TestSuiteHandler{
		testSuiteService: testSuiteService,
	}
}

// ListTestSuites handles listing all test suites
func (h *TestSuiteHandler) ListTestSuites(c *gin.Context) {
	suites, err := h.testSuiteService.ListTestSuites()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve test suites"})
		return
	}

	// Convert to response objects
	responses := make([]models.TestSuiteResponse, 0, len(suites))
	for _, suite := range suites {
		responses = append(responses, suite.ToResponse())
	}

	c.JSON(http.StatusOK, responses)
}

// CreateTestSuite handles creating a new test suite
func (h *TestSuiteHandler) CreateTestSuite(c *gin.Context) {
	var suiteCreate models.TestSuiteCreate
	if err := c.ShouldBindJSON(&suiteCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Create test suite from request data
	suite := &models.TestSuite{
		ProjectID:   suiteCreate.ProjectID,
		Name:        suiteCreate.Name,
		Description: suiteCreate.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := h.testSuiteService.CreateTestSuite(suite)
	if err != nil {
		if err == repository.ErrTestSuiteExists {
			c.JSON(http.StatusConflict, gin.H{"error": "Test suite with this name already exists in this project"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create test suite"})
		return
	}

	c.JSON(http.StatusCreated, suite.ToResponse())
}

// GetTestSuite handles retrieving a test suite by ID
func (h *TestSuiteHandler) GetTestSuite(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test suite ID"})
		return
	}

	suite, err := h.testSuiteService.GetTestSuiteByID(id)
	if err != nil {
		if err == repository.ErrTestSuiteNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Test suite not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve test suite"})
		return
	}

	c.JSON(http.StatusOK, suite.ToResponse())
}

// UpdateTestSuite handles updating a test suite
func (h *TestSuiteHandler) UpdateTestSuite(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test suite ID"})
		return
	}

	// Get existing test suite
	suite, err := h.testSuiteService.GetTestSuiteByID(id)
	if err != nil {
		if err == repository.ErrTestSuiteNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Test suite not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve test suite"})
		return
	}

	// Parse update data
	var suiteUpdate models.TestSuiteUpdate
	if err := c.ShouldBindJSON(&suiteUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Update fields if provided
	if suiteUpdate.Name != "" {
		suite.Name = suiteUpdate.Name
	}
	if suiteUpdate.Description != "" {
		suite.Description = suiteUpdate.Description
	}
	suite.UpdatedAt = time.Now()

	err = h.testSuiteService.UpdateTestSuite(suite)
	if err != nil {
		if err == repository.ErrTestSuiteExists {
			c.JSON(http.StatusConflict, gin.H{"error": "Test suite with this name already exists in this project"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update test suite"})
		return
	}

	c.JSON(http.StatusOK, suite.ToResponse())
}

// DeleteTestSuite handles deleting a test suite
func (h *TestSuiteHandler) DeleteTestSuite(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test suite ID"})
		return
	}

	err = h.testSuiteService.DeleteTestSuite(id)
	if err != nil {
		if err == repository.ErrTestSuiteNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Test suite not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete test suite"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Test suite deleted successfully"})
}

// ListTestSuitesByProject handles listing test suites by project
func (h *TestSuiteHandler) ListTestSuitesByProject(c *gin.Context) {
	projectID, err := strconv.ParseInt(c.Param("projectId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	suites, err := h.testSuiteService.ListTestSuitesByProject(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve test suites"})
		return
	}

	// Convert to response objects
	responses := make([]models.TestSuiteResponse, 0, len(suites))
	for _, suite := range suites {
		responses = append(responses, suite.ToResponse())
	}

	c.JSON(http.StatusOK, responses)
}
