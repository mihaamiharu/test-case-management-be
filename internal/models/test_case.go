package models

import (
	"time"
)

// TestCaseStatus represents the status of a test case
type TestCaseStatus string

const (
	StatusDraft      TestCaseStatus = "draft"
	StatusActive     TestCaseStatus = "active"
	StatusDeprecated TestCaseStatus = "deprecated"
)

// TestCasePriority represents the priority of a test case
type TestCasePriority string

const (
	PriorityLow    TestCasePriority = "low"
	PriorityMedium TestCasePriority = "medium"
	PriorityHigh   TestCasePriority = "high"
)

// StepType represents the type of a test step in Gherkin syntax
type StepType string

const (
	StepTypeGiven StepType = "given"
	StepTypeWhen  StepType = "when"
	StepTypeThen  StepType = "then"
	StepTypeAnd   StepType = "and"
	StepTypeBut   StepType = "but"
)

// TestCase represents a test case in the system
type TestCase struct {
	ID            int64            `json:"id"`
	ProjectID     int64            `json:"project_id"`
	SuiteID       int64            `json:"suite_id"`
	Title         string           `json:"title"`
	Description   string           `json:"description"`
	Preconditions string           `json:"preconditions"`
	Status        TestCaseStatus   `json:"status"`
	Priority      TestCasePriority `json:"priority"`
	CreatedBy     int64            `json:"created_by"`
	UpdatedBy     int64            `json:"updated_by"`
	Version       int              `json:"version"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
	Steps         []*TestStep      `json:"steps,omitempty"`
	Tags          []*Tag           `json:"tags,omitempty"`
}

// TestStep represents a step in a test case
type TestStep struct {
	ID             int64             `json:"id"`
	TestCaseID     int64             `json:"test_case_id"`
	StepNumber     int               `json:"step_number"`
	StepType       StepType          `json:"step_type"`
	Description    string            `json:"description"`
	ExpectedResult string            `json:"expected_result"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
	Notes          []*StepNote       `json:"notes,omitempty"`
	Attachments    []*StepAttachment `json:"attachments,omitempty"`
}

// StepNote represents a note attached to a test step
type StepNote struct {
	ID        int64     `json:"id"`
	StepID    int64     `json:"step_id"`
	Content   string    `json:"content"`
	CreatedBy int64     `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

// StepAttachment represents an attachment for a test step
type StepAttachment struct {
	ID        int64     `json:"id"`
	StepID    int64     `json:"step_id"`
	FileName  string    `json:"file_name"`
	FilePath  string    `json:"file_path"`
	FileType  string    `json:"file_type"`
	FileSize  int64     `json:"file_size"`
	CreatedBy int64     `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

// TestCaseHistory represents a historical version of a test case
type TestCaseHistory struct {
	ID            int64            `json:"id"`
	TestCaseID    int64            `json:"test_case_id"`
	Title         string           `json:"title"`
	Description   string           `json:"description"`
	Preconditions string           `json:"preconditions"`
	Status        TestCaseStatus   `json:"status"`
	Priority      TestCasePriority `json:"priority"`
	Version       int              `json:"version"`
	ChangedBy     int64            `json:"changed_by"`
	ChangeSummary string           `json:"change_summary"`
	CreatedAt     time.Time        `json:"created_at"`
}

// TestCaseCreate represents data needed to create a new test case
type TestCaseCreate struct {
	ProjectID     int64             `json:"project_id" binding:"required"`
	SuiteID       int64             `json:"suite_id" binding:"required"`
	Title         string            `json:"title" binding:"required"`
	Description   string            `json:"description"`
	Preconditions string            `json:"preconditions"`
	Status        TestCaseStatus    `json:"status" binding:"required,oneof=draft active deprecated"`
	Priority      TestCasePriority  `json:"priority" binding:"required,oneof=low medium high"`
	Steps         []*TestStepCreate `json:"steps"`
	Tags          []string          `json:"tags"`
}

// TestStepCreate represents data needed to create a new test step
type TestStepCreate struct {
	StepType       StepType `json:"step_type" binding:"required,oneof=given when then and but"`
	Description    string   `json:"description" binding:"required"`
	ExpectedResult string   `json:"expected_result"`
}

// TestCaseUpdate represents data needed to update a test case
type TestCaseUpdate struct {
	Title         string            `json:"title"`
	Description   string            `json:"description"`
	Preconditions string            `json:"preconditions"`
	Status        TestCaseStatus    `json:"status" binding:"omitempty,oneof=draft active deprecated"`
	Priority      TestCasePriority  `json:"priority" binding:"omitempty,oneof=low medium high"`
	Steps         []*TestStepCreate `json:"steps"`
	Tags          []string          `json:"tags"`
}

// StepNoteCreate represents data needed to create a new step note
type StepNoteCreate struct {
	Content string `json:"content" binding:"required"`
}

// TestCaseResponse represents the test case data to be returned in API responses
type TestCaseResponse struct {
	ID            int64               `json:"id"`
	ProjectID     int64               `json:"project_id"`
	SuiteID       int64               `json:"suite_id"`
	Title         string              `json:"title"`
	Description   string              `json:"description"`
	Preconditions string              `json:"preconditions"`
	Status        TestCaseStatus      `json:"status"`
	Priority      TestCasePriority    `json:"priority"`
	CreatedBy     int64               `json:"created_by"`
	UpdatedBy     int64               `json:"updated_by"`
	Version       int                 `json:"version"`
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
	Steps         []*TestStepResponse `json:"steps,omitempty"`
	Tags          []*TagResponse      `json:"tags,omitempty"`
}

// TestStepResponse represents the test step data to be returned in API responses
type TestStepResponse struct {
	ID             int64                     `json:"id"`
	StepNumber     int                       `json:"step_number"`
	StepType       StepType                  `json:"step_type"`
	Description    string                    `json:"description"`
	ExpectedResult string                    `json:"expected_result"`
	CreatedAt      time.Time                 `json:"created_at"`
	UpdatedAt      time.Time                 `json:"updated_at"`
	Notes          []*StepNoteResponse       `json:"notes,omitempty"`
	Attachments    []*StepAttachmentResponse `json:"attachments,omitempty"`
}

// StepNoteResponse represents the step note data to be returned in API responses
type StepNoteResponse struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	CreatedBy int64     `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

// StepAttachmentResponse represents the step attachment data to be returned in API responses
type StepAttachmentResponse struct {
	ID        int64     `json:"id"`
	FileName  string    `json:"file_name"`
	FileType  string    `json:"file_type"`
	FileSize  int64     `json:"file_size"`
	CreatedBy int64     `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse converts a TestCase to TestCaseResponse
func (tc *TestCase) ToResponse() *TestCaseResponse {
	response := &TestCaseResponse{
		ID:            tc.ID,
		ProjectID:     tc.ProjectID,
		SuiteID:       tc.SuiteID,
		Title:         tc.Title,
		Description:   tc.Description,
		Preconditions: tc.Preconditions,
		Status:        tc.Status,
		Priority:      tc.Priority,
		CreatedBy:     tc.CreatedBy,
		UpdatedBy:     tc.UpdatedBy,
		Version:       tc.Version,
		CreatedAt:     tc.CreatedAt,
		UpdatedAt:     tc.UpdatedAt,
	}

	if tc.Steps != nil {
		response.Steps = make([]*TestStepResponse, len(tc.Steps))
		for i, step := range tc.Steps {
			response.Steps[i] = step.ToResponse()
		}
	}

	if tc.Tags != nil {
		response.Tags = make([]*TagResponse, len(tc.Tags))
		for i, tag := range tc.Tags {
			response.Tags[i] = tag.ToResponse()
		}
	}

	return response
}

// ToResponse converts a TestStep to TestStepResponse
func (ts *TestStep) ToResponse() *TestStepResponse {
	response := &TestStepResponse{
		ID:             ts.ID,
		StepNumber:     ts.StepNumber,
		StepType:       ts.StepType,
		Description:    ts.Description,
		ExpectedResult: ts.ExpectedResult,
		CreatedAt:      ts.CreatedAt,
		UpdatedAt:      ts.UpdatedAt,
	}

	if ts.Notes != nil {
		response.Notes = make([]*StepNoteResponse, len(ts.Notes))
		for i, note := range ts.Notes {
			response.Notes[i] = note.ToResponse()
		}
	}

	if ts.Attachments != nil {
		response.Attachments = make([]*StepAttachmentResponse, len(ts.Attachments))
		for i, attachment := range ts.Attachments {
			response.Attachments[i] = attachment.ToResponse()
		}
	}

	return response
}

// ToResponse converts a StepNote to StepNoteResponse
func (sn *StepNote) ToResponse() *StepNoteResponse {
	return &StepNoteResponse{
		ID:        sn.ID,
		Content:   sn.Content,
		CreatedBy: sn.CreatedBy,
		CreatedAt: sn.CreatedAt,
	}
}

// ToResponse converts a StepAttachment to StepAttachmentResponse
func (a *StepAttachment) ToResponse() *StepAttachmentResponse {
	return &StepAttachmentResponse{
		ID:        a.ID,
		FileName:  a.FileName,
		FileType:  a.FileType,
		FileSize:  a.FileSize,
		CreatedBy: a.CreatedBy,
		CreatedAt: a.CreatedAt,
	}
}
