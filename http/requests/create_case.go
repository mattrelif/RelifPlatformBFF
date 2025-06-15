package requests

import (
	"relif/platform-bff/entities"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type CreateCaseInitialNote struct {
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	NoteType    string   `json:"note_type"`
	IsImportant bool     `json:"is_important"`
	Tags        []string `json:"tags"`
}

type CreateCase struct {
	BeneficiaryID     string                 `json:"beneficiary_id"`
	AssignedToID      string                 `json:"assigned_to_id"`
	Title             string                 `json:"title"`
	Description       string                 `json:"description"`
	CaseType          string                 `json:"case_type"`
	Priority          string                 `json:"priority"`
	UrgencyLevel      string                 `json:"urgency_level"`
	DueDate           string                 `json:"due_date"` // ISO date string
	EstimatedDuration string                 `json:"estimated_duration"`
	BudgetAllocated   string                 `json:"budget_allocated"`
	Tags              []string               `json:"tags"`
	InitialNote       *CreateCaseInitialNote `json:"initial_note"`
}

func (req *CreateCase) Validate() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.BeneficiaryID, validation.Required, is.MongoID),
		validation.Field(&req.AssignedToID, validation.Required, is.MongoID),
		validation.Field(&req.Title, validation.Required),
		validation.Field(&req.Description, validation.Required),
		validation.Field(&req.CaseType, validation.Required, validation.In(
			"HOUSING", "LEGAL", "MEDICAL", "SUPPORT", "EDUCATION",
			"EMPLOYMENT", "FINANCIAL", "FAMILY_REUNIFICATION",
			"DOCUMENTATION", "MENTAL_HEALTH", "OTHER",
		)),
		validation.Field(&req.Priority, validation.Required, validation.In(
			"LOW", "MEDIUM", "HIGH", "URGENT",
		)),
		validation.Field(&req.UrgencyLevel, validation.In(
			"IMMEDIATE", "WITHIN_WEEK", "WITHIN_MONTH", "FLEXIBLE", "",
		)),
		validation.Field(&req.InitialNote, validation.By(func(value interface{}) error {
			if note, ok := value.(*CreateCaseInitialNote); ok && note != nil {
				return validation.ValidateStruct(note,
					validation.Field(&note.Title, validation.Required),
					validation.Field(&note.Content, validation.Required),
					validation.Field(&note.NoteType, validation.Required, validation.In(
						"CALL", "MEETING", "UPDATE", "APPOINTMENT", "OTHER",
					)),
				)
			}
			return nil
		})),
	)
}

// ValidateOrganizationBoundaries ensures beneficiary and assigned user belong to the organization
func (req *CreateCase) ValidateOrganizationBoundaries(organizationID string, beneficiary entities.Beneficiary, assignedUser entities.User) error {
	if beneficiary.CurrentOrganizationID != organizationID {
		return validation.NewError("beneficiary_id", "beneficiary must belong to your organization")
	}

	if assignedUser.OrganizationID != organizationID {
		return validation.NewError("assigned_to_id", "assigned user must belong to your organization")
	}

	return nil
}

func (req *CreateCase) ToEntity(organizationID string) entities.Case {
	var dueDate *time.Time
	if req.DueDate != "" {
		if parsed, err := time.Parse("2006-01-02", req.DueDate); err == nil {
			dueDate = &parsed
		}
	}

	return entities.Case{
		BeneficiaryID:     req.BeneficiaryID,
		AssignedToID:      req.AssignedToID,
		Title:             req.Title,
		Description:       req.Description,
		CaseType:          req.CaseType,
		Priority:          req.Priority,
		UrgencyLevel:      req.UrgencyLevel,
		DueDate:           dueDate,
		EstimatedDuration: req.EstimatedDuration,
		BudgetAllocated:   req.BudgetAllocated,
		Tags:              req.Tags,
		OrganizationID:    organizationID,
	}
}

func (req *CreateCase) ToInitialNoteEntity(caseID, creatorID string) *entities.CaseNote {
	if req.InitialNote == nil {
		return nil
	}

	return &entities.CaseNote{
		CaseID:      caseID,
		Title:       req.InitialNote.Title,
		Content:     req.InitialNote.Content,
		Tags:        req.InitialNote.Tags,
		NoteType:    req.InitialNote.NoteType,
		IsImportant: req.InitialNote.IsImportant,
		CreatedByID: creatorID,
	}
}
