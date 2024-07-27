package entities

import "time"

type JoinOrganizationRequest struct {
	ID             string
	UserID         string
	OrganizationID string
	CreatedAt      time.Time
	ExpiresAt      time.Time
}
