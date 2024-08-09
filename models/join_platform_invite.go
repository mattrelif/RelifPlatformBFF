package models

import (
	"relif/platform-bff/entities"
	"time"
)

type JoinPlatformInvite struct {
	InvitedEmail   string    `bson:"_id,omitempty"`
	Code           string    `bson:"code,omitempty"`
	OrganizationID string    `bson:"organization_id,omitempty"`
	InviterID      string    `bson:"inviter_id,omitempty"`
	CreatedAt      time.Time `bson:"created_at,omitempty"`
	ExpiresAt      time.Time `bson:"expires_at,omitempty"`
}

func (invite *JoinPlatformInvite) ToEntity() entities.JoinPlatformInvite {
	return entities.JoinPlatformInvite{
		InvitedEmail:   invite.InvitedEmail,
		Code:           invite.Code,
		OrganizationID: invite.OrganizationID,
		InviterID:      invite.InviterID,
		CreatedAt:      invite.CreatedAt,
		ExpiresAt:      invite.ExpiresAt,
	}
}

func NewJoinPlatformInvite(entity entities.JoinPlatformInvite) JoinPlatformInvite {
	return JoinPlatformInvite{
		InvitedEmail:   entity.InvitedEmail,
		Code:           entity.Code,
		OrganizationID: entity.OrganizationID,
		InviterID:      entity.InviterID,
		CreatedAt:      time.Now(),
		ExpiresAt:      time.Now().Add(4 * time.Hour),
	}
}
