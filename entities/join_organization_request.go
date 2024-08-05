package entities

import "time"

type JoinOrganizationRequest struct {
	ID             string
	UserID         string
	OrganizationID string
	Status         string
	AuditorID      string
	CreatedAt      time.Time
	AcceptedAt     time.Time
	RejectedAt     time.Time
	RejectReason   string
	ExpiresAt      *time.Time
}
