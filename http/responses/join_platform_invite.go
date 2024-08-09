package responses

import (
	"relif/platform-bff/entities"
	"time"
)

type JoinPlatformInvites []JoinPlatformInvite

type JoinPlatformInvite struct {
	InvitedEmail   string    `json:"invited_email"`
	OrganizationID string    `json:"organization_id"`
	InviterID      string    `json:"inviter_id"`
	CreatedAt      time.Time `json:"created_at"`
	ExpiresAt      time.Time `json:"expires_at"`
}

func NewJoinPlatformInvite(entity entities.JoinPlatformInvite) JoinPlatformInvite {
	return JoinPlatformInvite{
		InvitedEmail:   entity.InvitedEmail,
		OrganizationID: entity.OrganizationID,
		InviterID:      entity.InviterID,
		CreatedAt:      entity.CreatedAt,
		ExpiresAt:      entity.ExpiresAt,
	}
}

func NewJoinPlatformInvites(entityList []entities.JoinPlatformInvite) JoinPlatformInvites {
	res := make(JoinPlatformInvites, 0)

	for _, entity := range entityList {
		res = append(res, NewJoinPlatformInvite(entity))
	}

	return res
}
