package entities

import "time"

type OrganizationDataAccessGrant struct {
	ID                   string
	TargetOrganizationID string
	TargetOrganization   Organization
	OrganizationID       string
	AuditorID            string
	CreatedAt            time.Time
}
