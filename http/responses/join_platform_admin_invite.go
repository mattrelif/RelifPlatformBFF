package responses

import (
	"relif/platform-bff/entities"
	"time"
)

type JoinPlatformAdminInvites []JoinPlatformAdminInvite

type JoinPlatformAdminInvite struct {
	InvitedEmail string        `json:"invited_email"`
	InviterID    string        `json:"inviter_id"`
	Inviter      entities.User `json:"inviter"`
	CreatedAt    time.Time     `json:"created_at"`
	ExpiresAt    time.Time     `json:"expires_at"`
}

func NewJoinPlatformAdminInvite(entity entities.JoinPlatformAdminInvite) JoinPlatformAdminInvite {
	return JoinPlatformAdminInvite{
		InvitedEmail: entity.InvitedEmail,
		InviterID:    entity.InviterID,
		Inviter:      entity.Inviter,
		CreatedAt:    entity.CreatedAt,
		ExpiresAt:    entity.ExpiresAt,
	}
}

func NewJoinPlatformAdminInvites(entityList []entities.JoinPlatformAdminInvite) JoinPlatformAdminInvites {
	res := make(JoinPlatformAdminInvites, 0)

	for _, entity := range entityList {
		res = append(res, NewJoinPlatformAdminInvite(entity))
	}

	return res
}
