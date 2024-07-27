package entities

import "time"

type OrganizationDataAccessRequest struct {
	ID                      string
	RequesterID             string
	RequesterOrganizationID string
	TargetOrganizationID    string
	AuditorID               string
	Status                  string
	CreatedAt               time.Time
	AcceptedAt              time.Time
	RejectReason            string
	RejectedAt              time.Time
}
