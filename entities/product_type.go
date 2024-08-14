package entities

import "time"

type ProductType struct {
	ID             string
	Name           string
	Description    string
	Brand          string
	Category       string
	OrganizationID string
	Organization   Organization
	UnitType       string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
