package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"relif/bff/entities"
	"relif/bff/utils"
	"time"
)

type FindHousingRoom struct {
	ID                string    `bson:"_id,omitempty"`
	HousingID         string    `bson:"housing_id,omitempty"`
	Name              string    `bson:"name,omitempty"`
	Status            string    `bson:"status,omitempty"`
	TotalVacancies    int       `bson:"total_vacancies,omitempty"`
	OccupiedVacancies int       `bson:"occupied_vacancies,omitempty"`
	CreatedAt         time.Time `bson:"created_at,omitempty"`
	UpdatedAt         time.Time `bson:"updated_at,omitempty"`
}

func (room *FindHousingRoom) ToEntity() entities.HousingRoom {
	return entities.HousingRoom{
		ID:                room.ID,
		HousingID:         room.HousingID,
		Name:              room.Name,
		Status:            room.Status,
		TotalVacancies:    room.TotalVacancies,
		OccupiedVacancies: room.OccupiedVacancies,
		CreatedAt:         room.CreatedAt,
		UpdatedAt:         room.UpdatedAt,
	}
}

type HousingRoom struct {
	ID             string    `bson:"_id,omitempty"`
	HousingID      string    `bson:"housing_id,omitempty"`
	Name           string    `bson:"name,omitempty"`
	Status         string    `bson:"status,omitempty"`
	TotalVacancies int       `bson:"total_vacancies,omitempty"`
	CreatedAt      time.Time `bson:"created_at,omitempty"`
	UpdatedAt      time.Time `bson:"updated_at,omitempty"`
}

func (room *HousingRoom) ToEntity() entities.HousingRoom {
	return entities.HousingRoom{
		ID:             room.ID,
		HousingID:      room.HousingID,
		Name:           room.Name,
		Status:         room.Status,
		TotalVacancies: room.TotalVacancies,
		CreatedAt:      room.CreatedAt,
		UpdatedAt:      room.UpdatedAt,
	}
}

func NewHousingRoom(entity entities.HousingRoom) HousingRoom {
	return HousingRoom{
		ID:             primitive.NewObjectID().Hex(),
		Name:           entity.Name,
		HousingID:      entity.HousingID,
		Status:         utils.ActiveStatus,
		TotalVacancies: entity.TotalVacancies,
		CreatedAt:      time.Now(),
	}
}

func NewUpdatedHousingRoom(entity entities.HousingRoom) HousingRoom {
	return HousingRoom{
		Name:           entity.Name,
		Status:         entity.Status,
		TotalVacancies: entity.TotalVacancies,
		UpdatedAt:      time.Now(),
	}
}
