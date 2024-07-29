package models

import (
	"relif/bff/entities"
	"time"
)

type JoinPlatformInvite struct {
	InvitedEmail   string    `bson:"_id"`
	Code           string    `bson:"code"`
	OrganizationID string    `bson:"organization_id"`
	InviterID      string    `bson:"inviter_id"`
	CreatedAt      time.Time `bson:"created_at"`
	ExpiresAt      time.Time `bson:"expires_at"`
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
