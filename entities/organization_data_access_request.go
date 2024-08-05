package entities

import "time"

type OrganizationDataAccessRequest struct {
	ID                      string
	RequesterID             string
	Requester               User
	RequesterOrganizationID string
	RequesterOrganization   Organization
	TargetOrganizationID    string
	TargetOrganization      Organization
	AuditorID               string
	Auditor                 User
	Status                  string
	CreatedAt               time.Time
	AcceptedAt              time.Time
	RejectReason            string
	RejectedAt              time.Time
}
