package entities

import "time"

type BeneficiaryStatusAudit struct {
	ID             string
	BeneficiaryID  string
	Beneficiary    Beneficiary
	PreviousStatus string
	NewStatus      string
	ChangedBy      string
	ChangedByUser  User
	ChangedAt      time.Time
	Reason         string
	OrganizationID string
	Organization   Organization
}
