package service

import (
	"github.com/mihaamiharu/test-case-management-be/internal/models"
	"github.com/mihaamiharu/test-case-management-be/internal/repository"
)

// TestCaseService handles test case business logic
type TestCaseService struct {
	testCaseRepo repository.TestCaseRepositoryInterface
	tagRepo      repository.TagRepositoryInterface
}

// NewTestCaseService creates a new test case service
func NewTestCaseService(testCaseRepo repository.TestCaseRepositoryInterface, tagRepo repository.TagRepositoryInterface) *TestCaseService {
	return &TestCaseService{
		testCaseRepo: testCaseRepo,
		tagRepo:      tagRepo,
	}
}

// CreateTestCase creates a new test case with tags
func (s *TestCaseService) CreateTestCase(testCase *models.TestCase, tagIDs []int64) error {
	// Create the test case
	if err := s.testCaseRepo.Create(testCase); err != nil {
		return err
	}

	// Add tags if provided
	if len(tagIDs) > 0 {
		for _, tagID := range tagIDs {
			if err := s.tagRepo.AddTagToTestCase(testCase.ID, tagID); err != nil {
				return err
			}
		}
	}

	return nil
}

// GetTestCaseByID retrieves a test case by ID
func (s *TestCaseService) GetTestCaseByID(id int64) (*models.TestCase, error) {
	testCase, err := s.testCaseRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Get tags for the test case
	tags, err := s.tagRepo.GetTagsByTestCase(id)
	if err != nil {
		return nil, err
	}
	testCase.Tags = tags

	return testCase, nil
}

// UpdateTestCase updates a test case with tags
func (s *TestCaseService) UpdateTestCase(testCase *models.TestCase, tagIDs []int64) error {
	// Update the test case
	if err := s.testCaseRepo.Update(testCase); err != nil {
		return err
	}

	// Update tags if provided
	if tagIDs != nil {
		if err := s.tagRepo.UpdateTestCaseTags(testCase.ID, tagIDs); err != nil {
			return err
		}
	}

	return nil
}

// DeleteTestCase deletes a test case
func (s *TestCaseService) DeleteTestCase(id int64) error {
	return s.testCaseRepo.Delete(id)
}

// ListTestCasesByProject retrieves all test cases for a project
func (s *TestCaseService) ListTestCasesByProject(projectID int64) ([]*models.TestCase, error) {
	testCases, err := s.testCaseRepo.ListByProject(projectID)
	if err != nil {
		return nil, err
	}

	// Get tags for each test case
	for _, tc := range testCases {
		tags, err := s.tagRepo.GetTagsByTestCase(tc.ID)
		if err != nil {
			return nil, err
		}
		tc.Tags = tags
	}

	return testCases, nil
}

// ListTestCasesBySuite retrieves all test cases for a suite
func (s *TestCaseService) ListTestCasesBySuite(suiteID int64) ([]*models.TestCase, error) {
	testCases, err := s.testCaseRepo.ListBySuite(suiteID)
	if err != nil {
		return nil, err
	}

	// Get tags for each test case
	for _, tc := range testCases {
		tags, err := s.tagRepo.GetTagsByTestCase(tc.ID)
		if err != nil {
			return nil, err
		}
		tc.Tags = tags
	}

	return testCases, nil
}

// AddTestStep adds a step to a test case
func (s *TestCaseService) AddTestStep(testCaseID int64, step *models.TestStep) error {
	// Ensure the test case exists
	_, err := s.testCaseRepo.GetByID(testCaseID)
	if err != nil {
		return err
	}

	// Get the current highest step number
	steps, err := s.testCaseRepo.GetSteps(testCaseID)
	if err != nil {
		return err
	}

	// Set the step number to be the next in sequence
	if len(steps) > 0 {
		highestStepNumber := 0
		for _, existingStep := range steps {
			if existingStep.StepNumber > highestStepNumber {
				highestStepNumber = existingStep.StepNumber
			}
		}
		step.StepNumber = highestStepNumber + 1
	} else {
		step.StepNumber = 1
	}

	return s.testCaseRepo.CreateStep(step)
}

// UpdateTestStep updates a test step
func (s *TestCaseService) UpdateTestStep(stepID int64, step *models.TestStep) error {
	// Get the existing step to verify it exists and get its test case ID
	existingStep, err := s.testCaseRepo.GetStepByID(stepID)
	if err != nil {
		return err
	}

	// Preserve the step number and test case ID
	step.StepNumber = existingStep.StepNumber
	step.TestCaseID = existingStep.TestCaseID

	return s.testCaseRepo.UpdateStep(step)
}

// DeleteTestStep deletes a test step
func (s *TestCaseService) DeleteTestStep(stepID int64) error {
	return s.testCaseRepo.DeleteStep(stepID)
}

// AddStepNote adds a note to a test step
func (s *TestCaseService) AddStepNote(stepID int64, note *models.StepNote) error {
	// Verify the step exists
	_, err := s.testCaseRepo.GetStepByID(stepID)
	if err != nil {
		return err
	}

	return s.testCaseRepo.CreateStepNote(note)
}

// DeleteStepNote deletes a step note
func (s *TestCaseService) DeleteStepNote(noteID int64) error {
	return s.testCaseRepo.DeleteStepNote(noteID)
}

// AddStepAttachment adds an attachment to a test step
func (s *TestCaseService) AddStepAttachment(stepID int64, attachment *models.StepAttachment) error {
	// Verify the step exists
	_, err := s.testCaseRepo.GetStepByID(stepID)
	if err != nil {
		return err
	}

	return s.testCaseRepo.CreateStepAttachment(attachment)
}

// GetStepAttachment gets an attachment by ID
func (s *TestCaseService) GetStepAttachment(attachmentID int64) (*models.StepAttachment, error) {
	// We need to implement this in the repository
	return s.testCaseRepo.GetStepAttachmentByID(attachmentID)
}

// DeleteStepAttachment deletes a step attachment
func (s *TestCaseService) DeleteStepAttachment(attachmentID int64) error {
	return s.testCaseRepo.DeleteStepAttachment(attachmentID)
}

// GetStepByID gets a step by ID with all its data
func (s *TestCaseService) GetStepByID(stepID int64) (*models.TestStep, error) {
	return s.testCaseRepo.GetStepByID(stepID)
}
