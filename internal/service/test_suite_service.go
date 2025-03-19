package service

import (
	"github.com/mihaamiharu/test-case-management-be/internal/models"
	"github.com/mihaamiharu/test-case-management-be/internal/repository"
)

// TestSuiteService handles test suite business logic
type TestSuiteService struct {
	testSuiteRepo repository.TestSuiteRepositoryInterface
}

// NewTestSuiteService creates a new test suite service
func NewTestSuiteService(testSuiteRepo repository.TestSuiteRepositoryInterface) *TestSuiteService {
	return &TestSuiteService{
		testSuiteRepo: testSuiteRepo,
	}
}

// CreateTestSuite creates a new test suite
func (s *TestSuiteService) CreateTestSuite(suite *models.TestSuite) error {
	return s.testSuiteRepo.Create(suite)
}

// GetTestSuiteByID retrieves a test suite by ID
func (s *TestSuiteService) GetTestSuiteByID(id int64) (*models.TestSuite, error) {
	return s.testSuiteRepo.GetByID(id)
}

// UpdateTestSuite updates a test suite
func (s *TestSuiteService) UpdateTestSuite(suite *models.TestSuite) error {
	return s.testSuiteRepo.Update(suite)
}

// DeleteTestSuite deletes a test suite
func (s *TestSuiteService) DeleteTestSuite(id int64) error {
	return s.testSuiteRepo.Delete(id)
}

// ListTestSuitesByProject retrieves all test suites for a project
func (s *TestSuiteService) ListTestSuitesByProject(projectID int64) ([]*models.TestSuite, error) {
	return s.testSuiteRepo.ListByProject(projectID)
}

// ListTestSuites retrieves all test suites
func (s *TestSuiteService) ListTestSuites() ([]*models.TestSuite, error) {
	return s.testSuiteRepo.List()
}
