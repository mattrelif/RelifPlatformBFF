package entities

import "time"

type Case struct {
	ID                string
	CaseNumber        string
	Title             string
	Description       string
	Status            string // "OPEN", "IN_PROGRESS", "PENDING", "ON_HOLD", "CLOSED"
	Priority          string // "LOW", "MEDIUM", "HIGH", "URGENT"
	UrgencyLevel      string // "IMMEDIATE", "WITHIN_WEEK", "WITHIN_MONTH", "FLEXIBLE"
	CaseType          string // "HOUSING", "LEGAL", "MEDICAL", "SUPPORT", etc.
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
