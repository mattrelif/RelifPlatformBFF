package entities

import "time"

type JoinOrganizationInvite struct {
	ID             string
	UserID         string
	User           User
	OrganizationID string
	Organization   Organization
	CreatorID      string
	Creator        User
	Status         string
	AcceptedAt     time.Time
	RejectedAt     time.Time
	RejectReason   string
	CreatedAt      time.Time
	ExpiresAt      *time.Time
}
