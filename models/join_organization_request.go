package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"relif/bff/entities"
	"relif/bff/utils"
	"time"
)

type JoinOrganizationRequest struct {
	ID             string     `bson:"_id,omitempty"`
	UserID         string     `bson:"user_id,omitempty"`
	OrganizationID string     `bson:"organization_id,omitempty"`
	Status         string     `bson:"status,omitempty"`
	AuditorID      string     `bson:"auditor_id,omitempty"`
	CreatedAt      time.Time  `bson:"created_at,omitempty"`
	AcceptedAt     time.Time  `bson:"accepted_at,omitempty"`
	RejectedAt     time.Time  `bson:"rejected_at,omitempty"`
	RejectReason   string     `bson:"reject_reason,omitempty"`
	ExpiresAt      *time.Time `bson:"expires_at,omitempty"`
}

func (request *JoinOrganizationRequest) ToEntity() entities.JoinOrganizationRequest {
	return entities.JoinOrganizationRequest{
		ID:             request.ID,
		UserID:         request.UserID,
		OrganizationID: request.OrganizationID,
		Status:         request.Status,
		AuditorID:      request.AuditorID,
		CreatedAt:      request.CreatedAt,
		AcceptedAt:     request.AcceptedAt,
		RejectedAt:     request.RejectedAt,
		RejectReason:   request.RejectReason,
		ExpiresAt:      request.ExpiresAt,
	}
}

func NewJoinOrganizationRequest(entity entities.JoinOrganizationRequest) JoinOrganizationRequest {
	expiresAt := time.Now().Add(4 * time.Hour)

	return JoinOrganizationRequest{
		ID:             primitive.NewObjectID().Hex(),
		UserID:         entity.UserID,
		OrganizationID: entity.OrganizationID,
		Status:         utils.PendingStatus,
		CreatedAt:      time.Now(),
		ExpiresAt:      &expiresAt,
	}
}

func NewUpdatedJoinOrganizationRequest(entity entities.JoinOrganizationRequest) JoinOrganizationRequest {
	return JoinOrganizationRequest{
		Status:       entity.Status,
		AuditorID:    entity.AuditorID,
		AcceptedAt:   entity.AcceptedAt,
		RejectedAt:   entity.RejectedAt,
		RejectReason: entity.RejectReason,
		ExpiresAt:    nil,
	}
}
