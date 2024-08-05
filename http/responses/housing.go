package responses

import (
	"relif/bff/entities"
	"time"
)

type Housings []Housing

type Housing struct {
	ID                string    `json:"id"`
	OrganizationID    string    `json:"organization_id"`
	Name              string    `json:"name"`
	Status            string    `json:"status"`
	TotalVacancies    int       `json:"total_vacancies"`
	OccupiedVacancies int       `json:"occupied_vacancies"`
	Address           Address   `json:"address"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func NewHousing(housing entities.Housing) Housing {
	return Housing{
		ID:                housing.ID,
		OrganizationID:    housing.OrganizationID,
		Name:              housing.Name,
		Status:            housing.Status,
		TotalVacancies:    housing.TotalVacancies,
		OccupiedVacancies: housing.OccupiedVacancies,
		Address:           NewAddress(housing.Address),
		CreatedAt:         housing.CreatedAt,
		UpdatedAt:         housing.UpdatedAt,
	}
}

func NewHousings(entityList []entities.Housing) Housings {
	res := make(Housings, 0)

	for _, entity := range entityList {
		res = append(res, NewHousing(entity))
	}

	return res
}
