package models

import (
	"relif/platform-bff/entities"
	"time"
)

type FindJoinPlatformAdminInvite struct {
	InvitedEmail string    `bson:"_id,omitempty"`
	Code         string    `bson:"code,omitempty"`
	InviterID    string    `bson:"inviter_id,omitempty"`
	Inviter      User      `bson:"inviter,omitempty"`
	CreatedAt    time.Time `bson:"created_at,omitempty"`
	ExpiresAt    time.Time `bson:"expires_at,omitempty"`
}

func (invite *FindJoinPlatformAdminInvite) ToEntity() entities.JoinPlatformAdminInvite {
	return entities.JoinPlatformAdminInvite{
		InvitedEmail: invite.InvitedEmail,
		Code:         invite.Code,
		InviterID:    invite.InviterID,
		Inviter:      invite.Inviter.ToEntity(),
		CreatedAt:    invite.CreatedAt,
		ExpiresAt:    invite.ExpiresAt,
	}
}

type JoinPlatformAdminInvite struct {
	InvitedEmail string    `bson:"_id"`
	Code         string    `bson:"code"`
	InviterID    string    `bson:"inviter_id"`
	CreatedAt    time.Time `bson:"created_at"`
	ExpiresAt    time.Time `bson:"expires_at"`
}

func (invite *JoinPlatformAdminInvite) ToEntity() entities.JoinPlatformAdminInvite {
	return entities.JoinPlatformAdminInvite{
		InvitedEmail: invite.InvitedEmail,
		Code:         invite.Code,
		InviterID:    invite.InviterID,
		CreatedAt:    invite.CreatedAt,
		ExpiresAt:    invite.ExpiresAt,
	}
}

func NewJoinPlatformAdminInvite(entity entities.JoinPlatformAdminInvite) JoinPlatformAdminInvite {
	return JoinPlatformAdminInvite{
		InvitedEmail: entity.InvitedEmail,
		Code:         entity.Code,
		InviterID:    entity.InviterID,
		CreatedAt:    time.Now(),
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}
}
