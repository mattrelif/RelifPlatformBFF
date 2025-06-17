package models

import (
	"fmt"
	"relif/platform-bff/entities"
	"time"

	"relif/platform-bff/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Case struct {
	ID                string     `bson:"_id,omitempty"`
	CaseNumber        string     `bson:"case_number,omitempty"`
	Title             string     `bson:"title,omitempty"`
	Description       string     `bson:"description,omitempty"`
	Status            string     `bson:"status,omitempty"`
	Priority          string     `bson:"priority,omitempty"`
	UrgencyLevel      string     `bson:"urgency_level,omitempty"`
	CaseType          string     `bson:"case_type,omitempty"`     // DEPRECATED: Kept for backwards compatibility
	ServiceTypes      []string   `bson:"service_types,omitempty"` // New field: Array of humanitarian service types
	BeneficiaryID     string     `bson:"beneficiary_id,omitempty"`
	AssignedToID      string     `bson:"assigned_to_id,omitempty"`
	DueDate           *time.Time `bson:"due_date,omitempty"`
	EstimatedDuration string     `bson:"estimated_duration,omitempty"`
	BudgetAllocated   string     `bson:"budget_allocated,omitempty"`
	Tags              []string   `bson:"tags,omitempty"`
	NotesCount        int        `bson:"notes_count,omitempty"`
	DocumentsCount    int        `bson:"documents_count,omitempty"`
	LastActivity      time.Time  `bson:"last_activity,omitempty"`
	OrganizationID    string     `bson:"organization_id,omitempty"`
	CreatedAt         time.Time  `bson:"created_at,omitempty"`
	UpdatedAt         time.Time  `bson:"updated_at,omitempty"`
}

func (c *Case) ToEntity() entities.Case {
	// Migration logic: If ServiceTypes is empty but CaseType exists, migrate it
	serviceTypes := c.ServiceTypes
	if len(serviceTypes) == 0 && c.CaseType != "" {
		serviceTypes = utils.MigrateCaseTypeToServiceTypes(c.CaseType)
	}

	return entities.Case{
		ID:                c.ID,
		CaseNumber:        c.CaseNumber,
		Title:             c.Title,
		Description:       c.Description,
		Status:            c.Status,
		Priority:          c.Priority,
		UrgencyLevel:      c.UrgencyLevel,
		CaseType:          c.CaseType,
		ServiceTypes:      serviceTypes,
		BeneficiaryID:     c.BeneficiaryID,
		AssignedToID:      c.AssignedToID,
		DueDate:           c.DueDate,
		EstimatedDuration: c.EstimatedDuration,
		BudgetAllocated:   c.BudgetAllocated,
		Tags:              c.Tags,
		NotesCount:        c.NotesCount,
		DocumentsCount:    c.DocumentsCount,
		LastActivity:      c.LastActivity,
		OrganizationID:    c.OrganizationID,
		CreatedAt:         c.CreatedAt,
		UpdatedAt:         c.UpdatedAt,
	}
}

func NewCase(entity entities.Case) Case {
	// Backwards compatibility: If entity has CaseType but no ServiceTypes, migrate it
	serviceTypes := entity.ServiceTypes
	if len(serviceTypes) == 0 && entity.CaseType != "" {
		serviceTypes = utils.MigrateCaseTypeToServiceTypes(entity.CaseType)
	}

	return Case{
		ID:                primitive.NewObjectID().Hex(),
		CaseNumber:        generateCaseNumber(),
		Title:             entity.Title,
		Description:       entity.Description,
		Status:            "PENDING",
		Priority:          entity.Priority,
		UrgencyLevel:      entity.UrgencyLevel,
		CaseType:          entity.CaseType,
		ServiceTypes:      serviceTypes,
		BeneficiaryID:     entity.BeneficiaryID,
		AssignedToID:      entity.AssignedToID,
		DueDate:           entity.DueDate,
		EstimatedDuration: entity.EstimatedDuration,
		BudgetAllocated:   entity.BudgetAllocated,
		Tags:              entity.Tags,
		NotesCount:        0,
		DocumentsCount:    0,
		LastActivity:      time.Now(),
		OrganizationID:    entity.OrganizationID,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
}

func NewCaseFromEntity(entity entities.Case) *Case {
	// Backwards compatibility: If entity has CaseType but no ServiceTypes, migrate it
	serviceTypes := entity.ServiceTypes
	if len(serviceTypes) == 0 && entity.CaseType != "" {
		serviceTypes = utils.MigrateCaseTypeToServiceTypes(entity.CaseType)
	}

	return &Case{
		ID:                primitive.NewObjectID().Hex(),
		CaseNumber:        generateCaseNumber(),
		Title:             entity.Title,
		Description:       entity.Description,
		Status:            "PENDING",
		Priority:          entity.Priority,
		UrgencyLevel:      entity.UrgencyLevel,
		CaseType:          entity.CaseType,
		ServiceTypes:      serviceTypes,
		BeneficiaryID:     entity.BeneficiaryID,
		AssignedToID:      entity.AssignedToID,
		DueDate:           entity.DueDate,
		EstimatedDuration: entity.EstimatedDuration,
		BudgetAllocated:   entity.BudgetAllocated,
		Tags:              entity.Tags,
		NotesCount:        0,
		DocumentsCount:    0,
		LastActivity:      time.Now(),
		OrganizationID:    entity.OrganizationID,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
}

func NewUpdatedCase(entity entities.Case) Case {
	// Backwards compatibility: If entity has CaseType but no ServiceTypes, migrate it
	serviceTypes := entity.ServiceTypes
	if len(serviceTypes) == 0 && entity.CaseType != "" {
		serviceTypes = utils.MigrateCaseTypeToServiceTypes(entity.CaseType)
	}

	return Case{
		Title:             entity.Title,
		Description:       entity.Description,
		Status:            entity.Status,
		Priority:          entity.Priority,
		UrgencyLevel:      entity.UrgencyLevel,
		CaseType:          entity.CaseType,
		ServiceTypes:      serviceTypes,
		AssignedToID:      entity.AssignedToID,
		DueDate:           entity.DueDate,
		EstimatedDuration: entity.EstimatedDuration,
		BudgetAllocated:   entity.BudgetAllocated,
		Tags:              entity.Tags,
		LastActivity:      time.Now(),
		UpdatedAt:         time.Now(),
		// NOTE: BeneficiaryID is intentionally excluded - beneficiaries cannot be changed after case creation
	}
}

// Generate case number in format CASE-YYYY-NNNN
func generateCaseNumber() string {
	year := time.Now().Year()
	// In production, this should query the database for the next sequence number
	// For now, using timestamp-based approach
	seq := time.Now().Unix() % 10000
	return fmt.Sprintf("CASE-%d-%04d", year, seq)
}
