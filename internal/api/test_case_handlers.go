package api

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mihaamiharu/test-case-management-be/internal/models"
	"github.com/mihaamiharu/test-case-management-be/internal/repository"
	"github.com/mihaamiharu/test-case-management-be/internal/service"
)

type TestCaseHandler struct {
	testCaseService *service.TestCaseService
}

func NewTestCaseHandler(testCaseService *service.TestCaseService) *TestCaseHandler {
	return &TestCaseHandler{
		testCaseService: testCaseService,
	}
}

// CreateTestCase handles the creation of a new test case
func (h *TestCaseHandler) CreateTestCase(c *gin.Context) {
	var testCaseCreate models.TestCaseCreate
	if err := c.ShouldBindJSON(&testCaseCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Create test case from request data
	testCase := &models.TestCase{
		ProjectID:     testCaseCreate.ProjectID,
		SuiteID:       testCaseCreate.SuiteID,
		Title:         testCaseCreate.Title,
		Description:   testCaseCreate.Description,
		Preconditions: testCaseCreate.Preconditions,
		Status:        testCaseCreate.Status,
		Priority:      testCaseCreate.Priority,
		CreatedBy:     userID.(int64),
		UpdatedBy:     userID.(int64),
	}

	// Convert step creates to steps
	if len(testCaseCreate.Steps) > 0 {
		steps := make([]*models.TestStep, len(testCaseCreate.Steps))
		for i, stepCreate := range testCaseCreate.Steps {
			steps[i] = &models.TestStep{
				StepType:       stepCreate.StepType,
				Description:    stepCreate.Description,
				ExpectedResult: stepCreate.ExpectedResult,
			}
		}
		testCase.Steps = steps
	}

	// Convert tag names to tag IDs
	var tagIDs []int64
	if len(testCaseCreate.Tags) > 0 {
		// This would typically involve looking up tag IDs by name
		// For now, we'll just use empty tag IDs
		tagIDs = []int64{}
	}

	// Create test case with tags
	err := h.testCaseService.CreateTestCase(testCase, tagIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, testCase.ToResponse())
}

// GetTestCase handles retrieving a test case by ID
func (h *TestCaseHandler) GetTestCase(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid test case ID"})
		return
	}

	testCase, err := h.testCaseService.GetTestCaseByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrTestCaseNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "test case not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, testCase.ToResponse())
}

// UpdateTestCase handles updating an existing test case
func (h *TestCaseHandler) UpdateTestCase(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid test case ID"})
		return
	}

	// Get existing test case
	testCase, err := h.testCaseService.GetTestCaseByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrTestCaseNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "test case not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Parse update data
	var testCaseUpdate models.TestCaseUpdate
	if err := c.ShouldBindJSON(&testCaseUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Update fields if provided
	if testCaseUpdate.Title != "" {
		testCase.Title = testCaseUpdate.Title
	}
	if testCaseUpdate.Description != "" {
		testCase.Description = testCaseUpdate.Description
	}
	if testCaseUpdate.Preconditions != "" {
		testCase.Preconditions = testCaseUpdate.Preconditions
	}
	if testCaseUpdate.Status != "" {
		testCase.Status = testCaseUpdate.Status
	}
	if testCaseUpdate.Priority != "" {
		testCase.Priority = testCaseUpdate.Priority
	}
	testCase.UpdatedBy = userID.(int64)

	// Update steps if provided
	if testCaseUpdate.Steps != nil {
		steps := make([]*models.TestStep, len(testCaseUpdate.Steps))
		for i, stepCreate := range testCaseUpdate.Steps {
			steps[i] = &models.TestStep{
				StepType:       stepCreate.StepType,
				Description:    stepCreate.Description,
				ExpectedResult: stepCreate.ExpectedResult,
			}
		}
		testCase.Steps = steps
	}

	// Convert tag names to tag IDs
	var tagIDs []int64
	if testCaseUpdate.Tags != nil {
		// This would typically involve looking up tag IDs by name
		// For now, we'll just use empty tag IDs
		tagIDs = []int64{}
	}

	// Update test case with tags
	err = h.testCaseService.UpdateTestCase(testCase, tagIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, testCase.ToResponse())
}

// DeleteTestCase handles deleting a test case
func (h *TestCaseHandler) DeleteTestCase(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid test case ID"})
		return
	}

	err = h.testCaseService.DeleteTestCase(id)
	if err != nil {
		if errors.Is(err, repository.ErrTestCaseNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "test case not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListTestCasesByProject handles listing test cases by project
func (h *TestCaseHandler) ListTestCasesByProject(c *gin.Context) {
	projectID, err := strconv.ParseInt(c.Param("projectId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project ID"})
		return
	}

	testCases, err := h.testCaseService.ListTestCasesByProject(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to response objects
	response := make([]*models.TestCaseResponse, len(testCases))
	for i, tc := range testCases {
		response[i] = tc.ToResponse()
	}

	// Always return an array (empty if no results)
	if response == nil {
		response = []*models.TestCaseResponse{}
	}

	c.JSON(http.StatusOK, response)
}

// ListTestCasesBySuite handles listing test cases by suite
func (h *TestCaseHandler) ListTestCasesBySuite(c *gin.Context) {
	suiteID, err := strconv.ParseInt(c.Param("suiteId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid suite ID"})
		return
	}

	testCases, err := h.testCaseService.ListTestCasesBySuite(suiteID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to response objects
	response := make([]*models.TestCaseResponse, len(testCases))
	for i, tc := range testCases {
		response[i] = tc.ToResponse()
	}

	// Always return an array (empty if no results)
	if response == nil {
		response = []*models.TestCaseResponse{}
	}

	c.JSON(http.StatusOK, response)
}

// AddTestStep handles adding a step to a test case
func (h *TestCaseHandler) AddTestStep(c *gin.Context) {
	testCaseID, err := strconv.ParseInt(c.Param("testCaseId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid test case ID"})
		return
	}

	var stepCreate models.TestStepCreate
	if err := c.ShouldBindJSON(&stepCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create step from request data
	step := &models.TestStep{
		TestCaseID:     testCaseID,
		StepNumber:     1, // Default to 1, would be determined by service
		StepType:       stepCreate.StepType,
		Description:    stepCreate.Description,
		ExpectedResult: stepCreate.ExpectedResult,
	}

	err = h.testCaseService.AddTestStep(testCaseID, step)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, step)
}

// UpdateTestStep handles updating a test step
func (h *TestCaseHandler) UpdateTestStep(c *gin.Context) {
	stepID, err := strconv.ParseInt(c.Param("stepId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid step ID"})
		return
	}

	var stepUpdate models.TestStepCreate
	if err := c.ShouldBindJSON(&stepUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create step from request data
	step := &models.TestStep{
		ID:             stepID,
		StepNumber:     1, // Default to 1, would be determined by service
		StepType:       stepUpdate.StepType,
		Description:    stepUpdate.Description,
		ExpectedResult: stepUpdate.ExpectedResult,
	}

	err = h.testCaseService.UpdateTestStep(stepID, step)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get the updated step with all its data (including attachments)
	updatedStep, err := h.testCaseService.GetStepByID(stepID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedStep.ToResponse())
}

// DeleteTestStep handles deleting a test step
func (h *TestCaseHandler) DeleteTestStep(c *gin.Context) {
	stepID, err := strconv.ParseInt(c.Param("stepId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid step ID"})
		return
	}

	err = h.testCaseService.DeleteTestStep(stepID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// AddStepNote handles adding a note to a test step
func (h *TestCaseHandler) AddStepNote(c *gin.Context) {
	stepID, err := strconv.ParseInt(c.Param("stepId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid step ID"})
		return
	}

	var noteCreate models.StepNoteCreate
	if err := c.ShouldBindJSON(&noteCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Create note from request data
	note := &models.StepNote{
		StepID:    stepID,
		Content:   noteCreate.Content,
		CreatedBy: userID.(int64),
	}

	err = h.testCaseService.AddStepNote(stepID, note)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get the updated step with all its data
	updatedStep, err := h.testCaseService.GetStepByID(stepID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, updatedStep.ToResponse())
}

// DeleteStepNote handles deleting a step note
func (h *TestCaseHandler) DeleteStepNote(c *gin.Context) {
	noteID, err := strconv.ParseInt(c.Param("noteId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid note ID"})
		return
	}

	err = h.testCaseService.DeleteStepNote(noteID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// UploadStepAttachment handles uploading an attachment to a test step
func (h *TestCaseHandler) UploadStepAttachment(c *gin.Context) {
	stepID, err := strconv.ParseInt(c.Param("stepId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid step ID"})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Get file from form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file uploaded"})
		return
	}
	defer file.Close()

	// Validate file size (e.g., max 10MB)
	if header.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file too large (max 10MB)"})
		return
	}

	// Get file extension and validate file type
	fileExt := filepath.Ext(header.Filename)
	fileType := c.Request.FormValue("file_type")
	if fileType == "" {
		// Try to determine file type from extension
		switch strings.ToLower(fileExt) {
		case ".jpg", ".jpeg", ".png", ".gif":
			fileType = "image"
		case ".pdf":
			fileType = "pdf"
		case ".doc", ".docx":
			fileType = "document"
		default:
			fileType = "other"
		}
	}

	// Create directory if it doesn't exist
	uploadDir := "./uploads/step_attachments"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create upload directory"})
		return
	}

	// Generate unique filename
	filename := fmt.Sprintf("%d_%s_%s", stepID, time.Now().Format("20060102150405"), header.Filename)
	filePath := filepath.Join(uploadDir, filename)

	// Save file
	out, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	// Create attachment record
	attachment := &models.StepAttachment{
		StepID:    stepID,
		FileName:  header.Filename,
		FilePath:  filePath,
		FileType:  fileType,
		FileSize:  header.Size,
		CreatedBy: userID.(int64),
	}

	err = h.testCaseService.AddStepAttachment(stepID, attachment)
	if err != nil {
		// Clean up file if database insert fails
		os.Remove(filePath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get the updated step with all its data
	updatedStep, err := h.testCaseService.GetStepByID(stepID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, updatedStep.ToResponse())
}

// DeleteStepAttachment handles deleting a step attachment
func (h *TestCaseHandler) DeleteStepAttachment(c *gin.Context) {
	attachmentID, err := strconv.ParseInt(c.Param("attachmentId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid attachment ID"})
		return
	}

	// Get the attachment to find the file path
	attachment, err := h.testCaseService.GetStepAttachment(attachmentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Delete the attachment record
	err = h.testCaseService.DeleteStepAttachment(attachmentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Delete the file
	if attachment != nil && attachment.FilePath != "" {
		os.Remove(attachment.FilePath)
	}

	c.Status(http.StatusNoContent)
}
