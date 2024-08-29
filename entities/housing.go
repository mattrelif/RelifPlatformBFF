package entities

import "time"

type Housing struct {
	ID                string
	OrganizationID    string
	Name              string
	Status            string
	TotalVacancies    int
	TotalRooms        int
	OccupiedVacancies int
	Address           Address
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
