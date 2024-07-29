package entities

import "time"

type JoinPlatformInvite struct {
	InvitedEmail   string
	Code           string
	OrganizationID string
	InviterID      string
	CreatedAt      time.Time
	ExpiresAt      time.Time
}
