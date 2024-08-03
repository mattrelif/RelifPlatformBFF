package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"relif/bff/entities"
	"relif/bff/utils"
	"time"
)

type UpdateOrganizationTypeRequest struct {
	ID             string    `bson:"_id,omitempty"`
	OrganizationID string    `bson:"organization_id,omitempty"`
	CreatorID      string    `bson:"creator_id,omitempty"`
	AuditorID      string    `bson:"auditor_id,omitempty"`
	Status         string    `json:"status,omitempty"`
	CreatedAt      time.Time `bson:"created_at,omitempty"`
	AcceptedAt     time.Time `bson:"accepted_at,omitempty"`
	RejectReason   string    `json:"reject_reason,omitempty"`
	RejectedAt     time.Time `json:"rejected_at,omitempty"`
}

func (request *UpdateOrganizationTypeRequest) ToEntity() entities.UpdateOrganizationTypeRequest {
	return entities.UpdateOrganizationTypeRequest{
		ID:             request.ID,
		OrganizationID: request.OrganizationID,
		CreatorID:      request.CreatorID,
		AuditorID:      request.AuditorID,
		Status:         request.Status,
		CreatedAt:      request.CreatedAt,
		RejectReason:   request.RejectReason,
		RejectedAt:     request.RejectedAt,
	}
}

func NewUpdateOrganizationTypeRequest(entity entities.UpdateOrganizationTypeRequest) UpdateOrganizationTypeRequest {
	return UpdateOrganizationTypeRequest{
		ID:             primitive.NewObjectID().Hex(),
		OrganizationID: entity.OrganizationID,
		CreatorID:      entity.CreatorID,
		Status:         utils.PendingStatus,
		CreatedAt:      time.Now(),
	}
}

func NewUpdatedUpdateOrganizationTypeRequest(entity entities.UpdateOrganizationTypeRequest) UpdateOrganizationTypeRequest {
	return UpdateOrganizationTypeRequest{
		AuditorID:    entity.AuditorID,
		AcceptedAt:   entity.AcceptedAt,
		RejectedAt:   entity.RejectedAt,
		RejectReason: entity.RejectReason,
		Status:       entity.Status,
	}
}
