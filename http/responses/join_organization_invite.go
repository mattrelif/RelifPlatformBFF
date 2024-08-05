package responses

import (
	"relif/bff/entities"
	"time"
)

type JoinOrganizationInvites []JoinOrganizationInvite

type JoinOrganizationInvite struct {
	ID             string       `json:"id"`
	UserID         string       `json:"user_id"`
	User           User         `json:"user"`
	OrganizationID string       `json:"organization_id"`
	Organization   Organization `json:"organization"`
	CreatorID      string       `json:"creator_id"`
	Creator        User         `json:"creator"`
	Status         string       `json:"status"`
	AcceptedAt     time.Time    `json:"accepted_at"`
	RejectedAt     time.Time    `json:"rejected_at"`
	RejectedReason string       `json:"rejected_reason"`
	CreatedAt      time.Time    `json:"created_at"`
	ExpiresAt      *time.Time   `json:"expires_at"`
}

func NewJoinOrganizationInvite(entity entities.JoinOrganizationInvite) JoinOrganizationInvite {
	return JoinOrganizationInvite{
		ID:             entity.ID,
		UserID:         entity.UserID,
		User:           NewUser(entity.User),
		OrganizationID: entity.OrganizationID,
		Organization:   NewOrganization(entity.Organization),
		CreatorID:      entity.CreatorID,
		Creator:        NewUser(entity.Creator),
		Status:         entity.Status,
		AcceptedAt:     entity.AcceptedAt,
		RejectedAt:     entity.RejectedAt,
		RejectedReason: entity.RejectReason,
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
