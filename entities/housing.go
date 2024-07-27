package entities

import "time"

type Housing struct {
	ID             string
	OrganizationID string
	Name           string
	Status         string
	Address        Address
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
