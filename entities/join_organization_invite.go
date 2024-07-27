package entities

import "time"

type JoinOrganizationInvite struct {
	ID             string
	UserID         string
	OrganizationID string
	CreatorID      string
	CreatedAt      time.Time
	ExpiresAt      time.Time
}
