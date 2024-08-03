package entities

import "time"

type UpdateOrganizationTypeRequest struct {
	ID             string
	OrganizationID string
	CreatorID      string
	AuditorID      string
	Status         string
	CreatedAt      time.Time
	AcceptedAt     time.Time
	RejectReason   string
	RejectedAt     time.Time
}
