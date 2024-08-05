package entities

import "time"

type UpdateOrganizationTypeRequest struct {
	ID             string
	OrganizationID string
	Organization   Organization
	CreatorID      string
	Creator        User
	AuditorID      string
	Auditor        User
	Status         string
	CreatedAt      time.Time
	AcceptedAt     time.Time
	RejectReason   string
	RejectedAt     time.Time
}
