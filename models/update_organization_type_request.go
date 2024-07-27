package models

import (
	"relif/bff/entities"
	"time"
)

type UpdateOrganizationTypeRequest struct {
	ID             string    `bson:"_id"`
	OrganizationID string    `bson:"organization_id"`
	CreatorID      string    `bson:"creator_id"`
	AuditorID      string    `bson:"auditor_id"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `bson:"created_at"`
	RejectReason   string    `json:"reject_reason"`
	RejectedAt     time.Time `json:"rejected_at"`
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
