package entities

import "time"

type Housing struct {
	ID                string
	OrganizationID    string
	Name              string
	Status            string
	TotalVacancies    int
	OccupiedVacancies int
	Address           Address
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
