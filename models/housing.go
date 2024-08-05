package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"relif/bff/entities"
	"relif/bff/utils"
	"time"
)

type FindHousing struct {
	ID                string    `bson:"_id,omitempty"`
	OrganizationID    string    `bson:"organization_id,omitempty"`
	Name              string    `bson:"name,omitempty"`
	Status            string    `bson:"status,omitempty"`
	TotalVacancies    int       `bson:"total_vacancies,omitempty"`
	OccupiedVacancies int       `bson:"occupied_vacancies,omitempty"`
	Address           Address   `bson:"address,omitempty"`
	CreatedAt         time.Time `bson:"created_at,omitempty"`
	UpdatedAt         time.Time `bson:"updated_at,omitempty"`
}

func (housing *FindHousing) ToEntity() entities.Housing {
	return entities.Housing{
		ID:                housing.ID,
		OrganizationID:    housing.OrganizationID,
		Name:              housing.Name,
		Status:            housing.Status,
		TotalVacancies:    housing.TotalVacancies,
		OccupiedVacancies: housing.OccupiedVacancies,
		Address:           housing.Address.ToEntity(),
		CreatedAt:         housing.CreatedAt,
		UpdatedAt:         housing.UpdatedAt,
	}
}

type Housing struct {
	ID             string    `bson:"_id,omitempty"`
	OrganizationID string    `bson:"organization_id,omitempty"`
	Name           string    `bson:"name,omitempty"`
	Status         string    `bson:"status,omitempty"`
	Address        Address   `bson:"address,omitempty"`
	CreatedAt      time.Time `bson:"created_at,omitempty"`
	UpdatedAt      time.Time `bson:"updated_at,omitempty"`
}

func (housing *Housing) ToEntity() entities.Housing {
	return entities.Housing{
		ID:             housing.ID,
		OrganizationID: housing.OrganizationID,
		Name:           housing.Name,
		Status:         housing.Status,
		Address:        housing.Address.ToEntity(),
		CreatedAt:      housing.CreatedAt,
		UpdatedAt:      housing.UpdatedAt,
	}
}

func NewHousing(entity entities.Housing) Housing {
	return Housing{
		ID:             primitive.NewObjectID().Hex(),
		OrganizationID: entity.OrganizationID,
		Name:           entity.Name,
		Status:         utils.ActiveStatus,
		Address:        NewAddress(entity.Address),
		CreatedAt:      time.Now(),
	}
}

func NewUpdatedHousing(entity entities.Housing) Housing {
	return Housing{
		Name:      entity.Name,
		Status:    entity.Status,
		Address:   NewAddress(entity.Address),
		UpdatedAt: time.Now(),
	}
}
