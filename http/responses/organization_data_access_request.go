package responses

import (
	"relif/bff/entities"
	"time"
)

type OrganizationDataAccessRequests []OrganizationDataAccessRequest

type OrganizationDataAccessRequest struct {
	ID                      string    `json:"id"`
	RequesterID             string    `json:"requester_id"`
	RequesterOrganizationID string    `json:"requester_organization_id"`
	TargetOrganizationID    string    `json:"target_organization_id"`
	AuditorID               string    `json:"auditor_id"`
	Status                  string    `json:"status"`
	CreatedAt               time.Time `json:"created_at"`
	AcceptedAt              time.Time `json:"accepted_at"`
	RejectReason            string    `json:"reject_reason"`
	RejectedAt              time.Time `json:"rejected_at"`
}

func NewOrganizationDataAccessRequest(entity entities.OrganizationDataAccessRequest) OrganizationDataAccessRequest {
	return OrganizationDataAccessRequest{
		ID:                      entity.ID,
		RequesterID:             entity.RequesterID,
		RequesterOrganizationID: entity.RequesterOrganizationID,
		TargetOrganizationID:    entity.TargetOrganizationID,
		AuditorID:               entity.AuditorID,
		Status:                  entity.Status,
		CreatedAt:               entity.CreatedAt,
		AcceptedAt:              entity.AcceptedAt,
		RejectReason:            entity.RejectReason,
		RejectedAt:              entity.RejectedAt,
	}
}

func NewNewOrganizationDataAccessRequests(entityList []entities.OrganizationDataAccessRequest) OrganizationDataAccessRequests {
	res := make(OrganizationDataAccessRequests, 0)

	for _, entity := range entityList {
		res = append(res, NewOrganizationDataAccessRequest(entity))
	}

	return res
}
