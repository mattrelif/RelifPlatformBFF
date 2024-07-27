package responses

import (
	"relif/bff/entities"
	"time"
)

type JoinOrganizationRequests []JoinOrganizationRequest

type JoinOrganizationRequest struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	OrganizationID string    `json:"organization_id"`
	CreatedAt      time.Time `json:"created_at"`
	ExpiresAt      time.Time `json:"expires_at"`
}

func NewJoinOrganizationRequest(entity entities.JoinOrganizationRequest) JoinOrganizationRequest {
	return JoinOrganizationRequest{
		ID:             entity.ID,
		UserID:         entity.UserID,
		OrganizationID: entity.OrganizationID,
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
