package models

import (
	"relif/bff/entities"
	"time"
)

type OrganizationDataAccess struct {
	ID                   string    `bson:"_id"`
	TargetOrganizationID string    `bson:"target_organization_id"`
	OrganizationID       string    `bson:"organization_id"`
	AuditorID            string    `bson:"auditor_id"`
	CreatedAt            time.Time `bson:"created_at"`
}

func (grant *OrganizationDataAccess) ToEntity() entities.OrganizationDataAccessGrant {
	return entities.OrganizationDataAccessGrant{
		ID:                   grant.ID,
		TargetOrganizationID: grant.TargetOrganizationID,
		OrganizationID:       grant.OrganizationID,
		AuditorID:            grant.AuditorID,
		CreatedAt:            grant.CreatedAt,
	}
}
