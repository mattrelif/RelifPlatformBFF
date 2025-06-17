package entities

import "time"

type Case struct {
	ID                string
	CaseNumber        string
	Title             string
	Description       string
	Status            string   // "IN_PROGRESS", "PENDING", "ON_HOLD", "CLOSED", "CANCELLED"
	Priority          string   // "LOW", "MEDIUM", "HIGH", "URGENT"
	UrgencyLevel      string   // "IMMEDIATE", "WITHIN_WEEK", "WITHIN_MONTH", "FLEXIBLE"
	CaseType          string   // DEPRECATED: Use ServiceTypes instead. Kept for backwards compatibility during migration
	ServiceTypes      []string // New field: Array of humanitarian service types
	BeneficiaryID     string
	Beneficiary       Beneficiary // Populated when needed
	AssignedToID      string
	AssignedTo        User // Populated when needed
	DueDate           *time.Time
	EstimatedDuration string
	BudgetAllocated   string
	Tags              []string
	NotesCount        int
	DocumentsCount    int
	LastActivity      time.Time
	OrganizationID    string
	Organization      Organization
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
