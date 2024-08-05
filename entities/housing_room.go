package entities

import "time"

type HousingRoom struct {
	ID                string
	HousingID         string
	Name              string
	Status            string
	TotalVacancies    int
	OccupiedVacancies int
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
