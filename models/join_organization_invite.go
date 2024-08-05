package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"relif/bff/entities"
	"relif/bff/utils"
	"time"
)

type JoinOrganizationInvite struct {
	ID             string     `bson:"_id,omitempty"`
	UserID         string     `bson:"user_id,omitempty"`
	OrganizationID string     `bson:"organization_id,omitempty"`
	CreatorID      string     `bson:"creator_id,omitempty"`
	Status         string     `bson:"status,omitempty"`
	AcceptedAt     time.Time  `bson:"accepted_at,omitempty"`
	RejectedAt     time.Time  `bson:"rejected_at,omitempty"`
	RejectReason   string     `bson:"reject_reason,omitempty"`
	CreatedAt      time.Time  `bson:"created_at,omitempty"`
	ExpiresAt      *time.Time `bson:"expires_at,omitempty"`
}

func (invite *JoinOrganizationInvite) ToEntity() entities.JoinOrganizationInvite {
	return entities.JoinOrganizationInvite{
		ID:             invite.ID,
		UserID:         invite.UserID,
		OrganizationID: invite.OrganizationID,
		CreatorID:      invite.CreatorID,
		Status:         invite.Status,
		AcceptedAt:     invite.AcceptedAt,
		RejectedAt:     invite.RejectedAt,
		CreatedAt:      invite.CreatedAt,
		ExpiresAt:      invite.ExpiresAt,
	}
}

func NewJoinOrganizationInvite(entity entities.JoinOrganizationInvite) JoinOrganizationInvite {
	expiresAt := time.Now().Add(time.Hour * 4)

	return JoinOrganizationInvite{
		ID:             primitive.NewObjectID().Hex(),
		UserID:         entity.UserID,
		OrganizationID: entity.OrganizationID,
		CreatorID:      entity.CreatorID,
		Status:         utils.PendingStatus,
		CreatedAt:      time.Now(),
		ExpiresAt:      &expiresAt,
	}
}

func NewUpdatedJoinOrganizationInvite(entity entities.JoinOrganizationInvite) JoinOrganizationInvite {
	return JoinOrganizationInvite{
		Status:       entity.Status,
		AcceptedAt:   entity.AcceptedAt,
		RejectedAt:   entity.RejectedAt,
		RejectReason: entity.RejectReason,
		ExpiresAt:    nil,
	}
}
