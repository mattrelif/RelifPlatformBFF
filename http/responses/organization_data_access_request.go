package responses

import (
	"relif/bff/entities"
	"time"
)

type OrganizationDataAccessRequests []OrganizationDataAccessRequest

type OrganizationDataAccessRequest struct {
	ID                      string       `json:"id"`
	RequesterID             string       `json:"requester_id"`
	Requester               User         `json:"requester"`
	RequesterOrganizationID string       `json:"requester_organization_id"`
	RequesterOrganization   Organization `json:"requester_organization"`
	TargetOrganizationID    string       `json:"target_organization_id"`
	TargetOrganization      Organization `json:"target_organization"`
	AuditorID               string       `json:"auditor_id"`
	Auditor                 User         `json:"auditor"`
	Status                  string       `json:"status"`
	CreatedAt               time.Time    `json:"created_at"`
	AcceptedAt              time.Time    `json:"accepted_at"`
	RejectReason            string       `json:"reject_reason"`
	RejectedAt              time.Time    `json:"rejected_at"`
}

func NewOrganizationDataAccessRequest(entity entities.OrganizationDataAccessRequest) OrganizationDataAccessRequest {
	return OrganizationDataAccessRequest{
		ID:                      entity.ID,
		RequesterID:             entity.RequesterID,
		Requester:               NewUser(entity.Requester),
		RequesterOrganizationID: entity.RequesterOrganizationID,
		RequesterOrganization:   NewOrganization(entity.RequesterOrganization),
		TargetOrganizationID:    entity.TargetOrganizationID,
		TargetOrganization:      NewOrganization(entity.TargetOrganization),
		AuditorID:               entity.AuditorID,
		Auditor:                 NewUser(entity.Auditor),
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
