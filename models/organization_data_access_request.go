package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"relif/platform-bff/entities"
	"relif/platform-bff/utils"
	"time"
)

type FindOrganizationDataAccessRequest struct {
	ID                      string       `bson:"_id,omitempty"`
	RequesterID             string       `bson:"requester_id,omitempty"`
	Request                 User         `bson:"requester,omitempty"`
	RequesterOrganizationID string       `bson:"requester_organization_id,omitempty"`
	RequesterOrganization   Organization `bson:"requester_organization,omitempty"`
	TargetOrganizationID    string       `bson:"target_organization_id,omitempty"`
	TargetOrganization      Organization `bson:"target_organization,omitempty"`
	AuditorID               string       `bson:"auditor_id,omitempty"`
	Auditor                 User         `bson:"auditor,omitempty"`
	Status                  string       `bson:"status,omitempty"`
	CreatedAt               time.Time    `bson:"created_at,omitempty"`
	AcceptedAt              time.Time    `bson:"accepted_at,omitempty"`
	RejectReason            string       `bson:"reject_reason,omitempty"`
	RejectedAt              time.Time    `bson:"rejected_at,omitempty"`
}

func (request *FindOrganizationDataAccessRequest) ToEntity() entities.OrganizationDataAccessRequest {
	return entities.OrganizationDataAccessRequest{
		ID:                      request.ID,
		RequesterID:             request.RequesterID,
		Requester:               request.Request.ToEntity(),
		RequesterOrganizationID: request.RequesterOrganizationID,
		RequesterOrganization:   request.RequesterOrganization.ToEntity(),
		TargetOrganizationID:    request.TargetOrganizationID,
		TargetOrganization:      request.TargetOrganization.ToEntity(),
		AuditorID:               request.AuditorID,
		Auditor:                 request.Auditor.ToEntity(),
		Status:                  request.Status,
		CreatedAt:               request.CreatedAt,
		AcceptedAt:              request.AcceptedAt,
		RejectedAt:              request.RejectedAt,
		RejectReason:            request.RejectReason,
	}
}

type OrganizationDataAccessRequest struct {
	ID                      string    `bson:"_id,omitempty"`
	RequesterID             string    `bson:"requester_id,omitempty"`
	RequesterOrganizationID string    `bson:"requester_organization_id,omitempty"`
	TargetOrganizationID    string    `bson:"target_organization_id,omitempty"`
	AuditorID               string    `bson:"auditor_id,omitempty"`
	Status                  string    `bson:"status,omitempty"`
	CreatedAt               time.Time `bson:"created_at,omitempty"`
	AcceptedAt              time.Time `bson:"accepted_at,omitempty"`
	RejectReason            string    `bson:"reject_reason,omitempty"`
	RejectedAt              time.Time `bson:"rejected_at,omitempty"`
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

func NewOrganizationDataAccessRequest(entity entities.OrganizationDataAccessRequest) OrganizationDataAccessRequest {
	return OrganizationDataAccessRequest{
		ID:                      primitive.NewObjectID().Hex(),
		RequesterID:             entity.RequesterID,
		RequesterOrganizationID: entity.RequesterOrganizationID,
		TargetOrganizationID:    entity.TargetOrganizationID,
		Status:                  utils.PendingStatus,
		CreatedAt:               time.Now(),
	}
}

func NewUpdatedOrganizationDataAccessRequest(entity entities.OrganizationDataAccessRequest) OrganizationDataAccessRequest {
	return OrganizationDataAccessRequest{
		AuditorID:    entity.AuditorID,
		Status:       entity.Status,
		AcceptedAt:   entity.AcceptedAt,
		RejectedAt:   entity.RejectedAt,
		RejectReason: entity.RejectReason,
	}
}
