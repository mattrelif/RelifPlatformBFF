package models

import (
	"relif/bff/entities"
	"time"
)

type OrganizationDataAccessRequest struct {
	ID                      string    `bson:"_id"`
	RequesterID             string    `bson:"requester_id"`
	RequesterOrganizationID string    `bson:"requester_organization_id"`
	TargetOrganizationID    string    `bson:"target_organization_id"`
	AuditorID               string    `bson:"auditor_id"`
	Status                  string    `bson:"status"`
	CreatedAt               time.Time `bson:"created_at"`
	AcceptedAt              time.Time `bson:"accepted_at"`
	RejectReason            string    `bson:"reject_reason"`
	RejectedAt              time.Time `bson:"rejected_at"`
}

func (request *OrganizationDataAccessRequest) ToEntity() entities.OrganizationDataAccessRequest {
	return entities.OrganizationDataAccessRequest{
		ID:                      request.ID,
		RequesterID:             request.RequesterID,
		RequesterOrganizationID: request.RequesterOrganizationID,
		TargetOrganizationID:    request.TargetOrganizationID,
		AuditorID:               request.AuditorID,
		Status:                  request.Status,
		CreatedAt:               request.CreatedAt,
		AcceptedAt:              request.AcceptedAt,
		RejectedAt:              request.RejectedAt,
		RejectReason:            request.RejectReason,
	}
}
