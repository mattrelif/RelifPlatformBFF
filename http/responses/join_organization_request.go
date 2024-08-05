package responses

import (
	"relif/bff/entities"
	"time"
)

type JoinOrganizationRequests []JoinOrganizationRequest

type JoinOrganizationRequest struct {
	ID             string     `json:"id"`
	UserID         string     `json:"user_id"`
	OrganizationID string     `json:"organization_id"`
	Status         string     `json:"status"`
	AuditorID      string     `json:"auditor_id"`
	AcceptedAt     time.Time  `json:"accepted_at"`
	RejectedAt     time.Time  `json:"rejected_at"`
	RejectedReason string     `json:"rejected_reason"`
	CreatedAt      time.Time  `json:"created_at"`
	ExpiresAt      *time.Time `json:"expires_at"`
}

func NewJoinOrganizationRequest(entity entities.JoinOrganizationRequest) JoinOrganizationRequest {
	return JoinOrganizationRequest{
		ID:             entity.ID,
		UserID:         entity.UserID,
		OrganizationID: entity.OrganizationID,
		Status:         entity.Status,
		AuditorID:      entity.AuditorID,
		AcceptedAt:     entity.AcceptedAt,
		RejectedAt:     entity.RejectedAt,
		RejectedReason: entity.RejectReason,
		CreatedAt:      entity.CreatedAt,
		ExpiresAt:      entity.ExpiresAt,
	}
}

func NewJoinOrganizationRequests(entityList []entities.JoinOrganizationRequest) JoinOrganizationRequests {
	var response JoinOrganizationRequests

	for _, entity := range entityList {
		response = append(response, NewJoinOrganizationRequest(entity))
	}

	return response
}
