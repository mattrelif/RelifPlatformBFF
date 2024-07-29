package responses

import (
	"relif/bff/entities"
	"time"
)

type HousingRooms []HousingRoom

type HousingRoom struct {
	ID                 string
	HousingID          string
	Name               string
	Status             string
	TotalVacancies     int
	AvailableVacancies int
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func NewHousingRoom(entity entities.HousingRoom) HousingRoom {
	return HousingRoom{
		ID:                 entity.ID,
		HousingID:          entity.HousingID,
		Name:               entity.Name,
		Status:             entity.Status,
		TotalVacancies:     entity.TotalVacancies,
		AvailableVacancies: entity.AvailableVacancies,
		CreatedAt:          entity.CreatedAt,
		UpdatedAt:          entity.UpdatedAt,
	}
}

func NewHousingRooms(entityList []entities.HousingRoom) HousingRooms {
	var res HousingRooms

	for _, entity := range entityList {
		res = append(res, NewHousingRoom(entity))
	}

	return res
}
