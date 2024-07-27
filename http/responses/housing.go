package responses

import (
	"relif/bff/entities"
	"time"
)

type Housings []Housing

type Housing struct {
	ID             string
	OrganizationID string
	Name           string
	Status         string
	Address        Address
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewHousing(housing entities.Housing) Housing {
	return Housing{
		ID:             housing.ID,
		OrganizationID: housing.OrganizationID,
		Name:           housing.Name,
		Status:         housing.Status,
		Address:        NewAddress(housing.Address),
		CreatedAt:      housing.CreatedAt,
		UpdatedAt:      housing.UpdatedAt,
	}
}

func NewHousings(entityList []entities.Housing) Housings {
	res := make(Housings, 0)

	for _, entity := range entityList {
		res = append(res, NewHousing(entity))
	}

	return res
}
