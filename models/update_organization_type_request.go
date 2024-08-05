package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"relif/bff/entities"
	"relif/bff/utils"
	"time"
)

type FindUpdateOrganizationTypeRequest struct {
	ID             string       `bson:"_id,omitempty"`
	OrganizationID string       `bson:"organization_id,omitempty"`
	Organization   Organization `bson:"organization,omitempty"`
	CreatorID      string       `bson:"creator_id,omitempty"`
	Creator        User         `bson:"creator,omitempty"`
	AuditorID      string       `bson:"auditor_id,omitempty"`
	Auditor        User         `bson:"auditor,omitempty"`
	Status         string       `json:"status,omitempty"`
	CreatedAt      time.Time    `bson:"created_at,omitempty"`
	AcceptedAt     time.Time    `bson:"accepted_at,omitempty"`
	RejectReason   string       `json:"reject_reason,omitempty"`
	RejectedAt     time.Time    `json:"rejected_at,omitempty"`
}

func (request *FindUpdateOrganizationTypeRequest) ToEntity() entities.UpdateOrganizationTypeRequest {
	return entities.UpdateOrganizationTypeRequest{
		ID:             request.ID,
		OrganizationID: request.OrganizationID,
		Organization:   request.Organization.ToEntity(),
		CreatorID:      request.CreatorID,
		Creator:        request.Creator.ToEntity(),
		AuditorID:      request.AuditorID,
		Auditor:        request.Auditor.ToEntity(),
		Status:         request.Status,
		CreatedAt:      request.CreatedAt,
		RejectReason:   request.RejectReason,
		RejectedAt:     request.RejectedAt,
	}
}

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
