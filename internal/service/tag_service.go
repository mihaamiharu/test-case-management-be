package service

import (
	"github.com/mihaamiharu/test-case-management-be/internal/models"
	"github.com/mihaamiharu/test-case-management-be/internal/repository"
)

// TagService handles tag business logic
type TagService struct {
	tagRepo repository.TagRepositoryInterface
}

// NewTagService creates a new tag service
func NewTagService(tagRepo repository.TagRepositoryInterface) *TagService {
	return &TagService{
		tagRepo: tagRepo,
	}
}

// CreateTag creates a new tag
func (s *TagService) CreateTag(tag *models.Tag) error {
	return s.tagRepo.Create(tag)
}

// GetTagByID retrieves a tag by ID
func (s *TagService) GetTagByID(id int64) (*models.Tag, error) {
	return s.tagRepo.GetByID(id)
}

// DeleteTag deletes a tag
func (s *TagService) DeleteTag(id int64) error {
	return s.tagRepo.Delete(id)
}

// ListTags retrieves all tags
func (s *TagService) ListTags() ([]*models.Tag, error) {
	return s.tagRepo.List()
}

// GetTagsByTestCase retrieves all tags for a test case
func (s *TagService) GetTagsByTestCase(testCaseID int64) ([]*models.Tag, error) {
	return s.tagRepo.GetTagsByTestCase(testCaseID)
}

// AddTagToTestCase adds a tag to a test case
func (s *TagService) AddTagToTestCase(testCaseID, tagID int64) error {
	return s.tagRepo.AddTagToTestCase(testCaseID, tagID)
}

// RemoveTagFromTestCase removes a tag from a test case
func (s *TagService) RemoveTagFromTestCase(testCaseID, tagID int64) error {
	return s.tagRepo.RemoveTagFromTestCase(testCaseID, tagID)
}

// UpdateTestCaseTags updates the tags for a test case
func (s *TagService) UpdateTestCaseTags(testCaseID int64, tagIDs []int64) error {
	return s.tagRepo.UpdateTestCaseTags(testCaseID, tagIDs)
}
