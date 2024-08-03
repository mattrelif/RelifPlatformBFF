package responses

import (
	"relif/bff/entities"
	"time"
)

type HousingRooms []HousingRoom

type HousingRoom struct {
	ID                 string    `json:"id"`
	HousingID          string    `json:"housing_id"`
	Name               string    `json:"name"`
	Status             string    `json:"status"`
	TotalVacancies     int       `json:"total_vacancies"`
	AvailableVacancies int       `json:"available_vacancies"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
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
