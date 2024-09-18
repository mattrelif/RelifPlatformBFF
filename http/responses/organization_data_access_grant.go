package responses

import (
	"relif/platform-bff/entities"
	"time"
)

type OrganizationDataAccessGrants []OrganizationDataAccessGrant

type OrganizationDataAccessGrant struct {
	ID                   string       `json:"id"`
	TargetOrganizationID string       `json:"target_organization_id"`
	TargetOrganization   Organization `json:"target_organization"`
	OrganizationID       string       `json:"organization_id"`
	AuditorID            string       `json:"auditor_id"`
	CreatedAt            time.Time    `json:"created_at"`
}

func NewOrganizationDataAccessGrant(entity entities.OrganizationDataAccessGrant) OrganizationDataAccessGrant {
	return OrganizationDataAccessGrant{
		ID:                   entity.ID,
		TargetOrganizationID: entity.TargetOrganizationID,
		TargetOrganization:   NewOrganization(entity.TargetOrganization),
		OrganizationID:       entity.OrganizationID,
		AuditorID:            entity.AuditorID,
		CreatedAt:            entity.CreatedAt,
	}
}

func NewOrganizationDataAccessGrants(entityList []entities.OrganizationDataAccessGrant) OrganizationDataAccessGrants {
	res := make(OrganizationDataAccessGrants, 0)

	for _, entity := range entityList {
		res = append(res, NewOrganizationDataAccessGrant(entity))
	}

	return res
}
