package entities

import "time"

type JoinOrganizationInvite struct {
	ID             string
	UserID         string
	OrganizationID string
	CreatorID      string
	Status         string
	AcceptedAt     time.Time
	RejectedAt     time.Time
	CreatedAt      time.Time
	ExpiresAt      *time.Time
}
