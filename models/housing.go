package models

import (
	"relif/bff/entities"
	"time"
)

type Housing struct {
	ID             string    `bson:"_id"`
	OrganizationID string    `bson:"organization_id"`
	Name           string    `bson:"name"`
	Status         string    `bson:"status"`
	Address        Address   `bson:"address"`
	CreatedAt      time.Time `bson:"created_at"`
	UpdatedAt      time.Time `bson:"updated_at"`
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
