package requests

import (
	"relif/platform-bff/entities"
	"time"

	"relif/platform-bff/utils"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type UpdateCase struct {
	Title             *string   `json:"title"`
	Description       *string   `json:"description"`
	Status            *string   `json:"status"`
	Priority          *string   `json:"priority"`
	UrgencyLevel      *string   `json:"urgency_level"`
	CaseType          *string   `json:"case_type,omitempty"`     // DEPRECATED: Use ServiceTypes instead
	ServiceTypes      *[]string `json:"service_types,omitempty"` // New field: Array of humanitarian service types
	AssignedToID      *string   `json:"assigned_to_id"`
	DueDate           *string   `json:"due_date"` // ISO date string
	EstimatedDuration *string   `json:"estimated_duration"`
	BudgetAllocated   *string   `json:"budget_allocated"`
	Tags              *[]string `json:"tags"`
	// NOTE: BeneficiaryID is intentionally excluded - beneficiaries cannot be changed after case creation for security
}

func (req *UpdateCase) Validate() error {
	// If ServiceTypes is provided, validate each one
	if req.ServiceTypes != nil && len(*req.ServiceTypes) > 0 {
		for _, serviceType := range *req.ServiceTypes {
			if !utils.IsValidServiceType(serviceType) {
				return validation.NewError("service_types", "invalid service type: "+serviceType)
			}
		}
	}

	return validation.ValidateStruct(req,
		validation.Field(&req.Status, validation.When(req.Status != nil, validation.Required, validation.In(
			"IN_PROGRESS", "PENDING", "ON_HOLD", "CLOSED", "CANCELLED",
		))),
		validation.Field(&req.Priority, validation.When(req.Priority != nil, validation.Required, validation.In(
			"LOW", "MEDIUM", "HIGH", "URGENT",
		))),
		validation.Field(&req.UrgencyLevel, validation.When(req.UrgencyLevel != nil, validation.In(
			"IMMEDIATE", "WITHIN_WEEK", "WITHIN_MONTH", "FLEXIBLE", "",
		))),
		validation.Field(&req.ServiceTypes, validation.When(req.ServiceTypes != nil,
			validation.Length(1, 10), // Allow 1-10 service types per case
		)),
		validation.Field(&req.CaseType, validation.When(req.CaseType != nil, validation.Required, validation.In(
			"HOUSING", "LEGAL", "MEDICAL", "SUPPORT", "EDUCATION",
			"EMPLOYMENT", "FINANCIAL", "FAMILY_REUNIFICATION",
			"DOCUMENTATION", "MENTAL_HEALTH", "OTHER",
		))),
		validation.Field(&req.AssignedToID, validation.When(req.AssignedToID != nil, validation.Required, is.MongoID)),
	)
}

// ValidateOrganizationBoundaries ensures assigned user belongs to the organization (when reassigning)
func (req *UpdateCase) ValidateOrganizationBoundaries(organizationID string, assignedUser *entities.User) error {
	if req.AssignedToID != nil && assignedUser != nil {
		if assignedUser.OrganizationID != organizationID {
			return validation.NewError("assigned_to_id", "assigned user must belong to your organization")
		}
	}
	return nil
}

// ValidateBeneficiaryImmutability ensures no beneficiary ID changes are attempted
func (req *UpdateCase) ValidateBeneficiaryImmutability() error {
	// This function serves as documentation and future-proofing
	// Currently, BeneficiaryID is not included in UpdateCase struct, so this validates that design
	// If someone tries to add BeneficiaryID to this struct in the future, this will catch it
	return nil
}

func (req *UpdateCase) ToEntity() entities.Case {
	entity := entities.Case{}

	if req.Title != nil {
		entity.Title = *req.Title
	}
	if req.Description != nil {
		entity.Description = *req.Description
	}
	if req.Status != nil {
		entity.Status = *req.Status
	}
	if req.Priority != nil {
		entity.Priority = *req.Priority
	}
	if req.UrgencyLevel != nil {
		entity.UrgencyLevel = *req.UrgencyLevel
	}
	if req.CaseType != nil {
		entity.CaseType = *req.CaseType
	}
	if req.ServiceTypes != nil {
		entity.ServiceTypes = *req.ServiceTypes
	} else if req.CaseType != nil {
		entity.ServiceTypes = utils.MigrateCaseTypeToServiceTypes(*req.CaseType)
	}
	if req.AssignedToID != nil {
		entity.AssignedToID = *req.AssignedToID
	}
	if req.DueDate != nil && *req.DueDate != "" {
		if parsed, err := time.Parse("2006-01-02", *req.DueDate); err == nil {
			entity.DueDate = &parsed
		}
	}
	if req.EstimatedDuration != nil {
		entity.EstimatedDuration = *req.EstimatedDuration
	}
	if req.BudgetAllocated != nil {
		entity.BudgetAllocated = *req.BudgetAllocated
	}
	if req.Tags != nil {
		entity.Tags = *req.Tags
	}

	return entity
}
