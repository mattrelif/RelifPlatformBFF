package entities

import "time"

type JoinOrganizationRequest struct {
	ID             string
	UserID         string
	User           User
	OrganizationID string
	Organization   Organization
	Status         string
	AuditorID      string
	Auditor        User
	CreatedAt      time.Time
	AcceptedAt     time.Time
	RejectedAt     time.Time
	RejectReason   string
	ExpiresAt      *time.Time
}
