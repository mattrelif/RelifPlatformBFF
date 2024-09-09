package entities

import "time"

type JoinPlatformAdminInvite struct {
	InvitedEmail string
	Code         string
	InviterID    string
	Inviter      User
	CreatedAt    time.Time
	ExpiresAt    time.Time
}
