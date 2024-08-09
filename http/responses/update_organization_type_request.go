package responses

import (
	"relif/platform-bff/entities"
	"time"
)

type UpdateOrganizationTypeRequests []UpdateOrganizationTypeRequest

type UpdateOrganizationTypeRequest struct {
	ID             string       `json:"id"`
	OrganizationID string       `json:"organization_id"`
	Organization   Organization `json:"organization"`
	CreatorID      string       `json:"creator_id"`
	Creator        User         `json:"creator"`
	AuditorID      string       `json:"auditor_id"`
	Auditor        User         `json:"auditor"`
	Status         string       `json:"status"`
	CreatedAt      time.Time    `json:"created_at"`
	AcceptedAt     time.Time    `json:"accepted_at"`
	RejectReason   string       `json:"reject_reason"`
	RejectedAt     time.Time    `json:"rejected_at"`
}

func NewUpdateOrganizationTypeRequest(request entities.UpdateOrganizationTypeRequest) UpdateOrganizationTypeRequest {
	return UpdateOrganizationTypeRequest{
		ID:             request.ID,
		OrganizationID: request.OrganizationID,
		Organization:   NewOrganization(request.Organization),
		CreatorID:      request.CreatorID,
		Creator:        NewUser(request.Creator),
		AuditorID:      request.AuditorID,
		Auditor:        NewUser(request.Auditor),
		Status:         request.Status,
		CreatedAt:      request.CreatedAt,
		AcceptedAt:     request.AcceptedAt,
		RejectReason:   request.RejectReason,
		RejectedAt:     request.RejectedAt,
	}
}

func NewUpdateOrganizationTypeRequests(requestList []entities.UpdateOrganizationTypeRequest) UpdateOrganizationTypeRequests {
	res := make(UpdateOrganizationTypeRequests, 0)

	for _, request := range requestList {
		res = append(res, NewUpdateOrganizationTypeRequest(request))
	}

	return res
}
