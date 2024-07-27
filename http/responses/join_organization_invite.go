package responses

import (
	"relif/bff/entities"
	"time"
)

type JoinOrganizationInvites []JoinOrganizationInvite

type JoinOrganizationInvite struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	OrganizationID string    `json:"organization_id"`
	CreatorID      string    `json:"creator_id"`
	CreatedAt      time.Time `json:"created_at"`
	ExpiresAt      time.Time `json:"expires_at"`
}

func NewJoinOrganizationInvite(entity entities.JoinOrganizationInvite) JoinOrganizationInvite {
	return JoinOrganizationInvite{
		ID:             entity.ID,
		UserID:         entity.UserID,
		OrganizationID: entity.OrganizationID,
		CreatorID:      entity.CreatorID,
		CreatedAt:      entity.CreatedAt,
		ExpiresAt:      entity.ExpiresAt,
	}
}

func NewJoinOrganizationInvites(entityList []entities.JoinOrganizationInvite) JoinOrganizationInvites {
	var response JoinOrganizationInvites

	for _, entity := range entityList {
		response = append(response, NewJoinOrganizationInvite(entity))
	}

	return response
}
