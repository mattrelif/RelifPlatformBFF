package models

import (
	"relif/bff/entities"
	"time"
)

type JoinOrganizationRequest struct {
	ID             string `bson:"_id"`
	UserID         string `bson:"user_id"`
	OrganizationID string `bson:"organization_id"`
	CreatedAt      time.Time
	ExpiresAt      time.Time
}

func (request *JoinOrganizationRequest) ToEntity() entities.JoinOrganizationRequest {
	return entities.JoinOrganizationRequest{
		ID:             request.ID,
		UserID:         request.UserID,
		OrganizationID: request.OrganizationID,
		CreatedAt:      request.CreatedAt,
		ExpiresAt:      request.ExpiresAt,
	}
}
