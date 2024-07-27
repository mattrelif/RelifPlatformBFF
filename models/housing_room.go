package models

import (
	"relif/bff/entities"
	"time"
)

type HousingRoom struct {
	ID                 string    `bson:"_id"`
	HousingID          string    `bson:"housing_id"`
	Name               string    `bson:"name"`
	Status             string    `bson:"status"`
	TotalVacancies     int       `bson:"total_vacancies"`
	AvailableVacancies int       `bson:"available_vacancies"`
	CreatedAt          time.Time `bson:"created_at"`
	UpdatedAt          time.Time `bson:"updated_at"`
}

func (room *HousingRoom) ToEntity() entities.HousingRoom {
	return entities.HousingRoom{
		ID:                 room.ID,
		HousingID:          room.HousingID,
		Name:               room.Name,
		Status:             room.Status,
		TotalVacancies:     room.TotalVacancies,
		AvailableVacancies: room.AvailableVacancies,
		CreatedAt:          room.CreatedAt,
		UpdatedAt:          room.UpdatedAt,
	}
}
