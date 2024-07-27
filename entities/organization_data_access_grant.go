package entities

import "time"

type OrganizationDataAccessGrant struct {
	ID                   string
	TargetOrganizationID string
	OrganizationID       string
	AuditorID            string
	CreatedAt            time.Time
}
