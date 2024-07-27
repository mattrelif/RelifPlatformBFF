package models

import (
	"relif/bff/entities"
	"time"
)

type JoinOrganizationInvite struct {
	ID             string `bson:"_id"`
	UserID         string `bson:"user_id"`
	OrganizationID string `bson:"organization_id"`
	CreatorID      string `bson:"creator_id"`
	CreatedAt      time.Time
	ExpiresAt      time.Time
}

func (invite *JoinOrganizationInvite) ToEntity() entities.JoinOrganizationInvite {
	return entities.JoinOrganizationInvite{
		ID:             invite.ID,
		UserID:         invite.UserID,
		OrganizationID: invite.OrganizationID,
		CreatorID:      invite.CreatorID,
		CreatedAt:      invite.CreatedAt,
		ExpiresAt:      invite.ExpiresAt,
	}
}
