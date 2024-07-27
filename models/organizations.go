package models

import (
	"relif/bff/entities"
	"time"
)

type Organization struct {
	ID          string    `bson:"_id"`
	Name        string    `bson:"name"`
	Description string    `bson:"description"`
	Address     Address   `bson:"address"`
	Type        string    `bson:"type"`
	CreatorID   string    `bson:"creator_id"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
}

func (organization *Organization) ToEntity() entities.Organization {
	return entities.Organization{
		ID:          organization.ID,
		Name:        organization.Name,
		Description: organization.Description,
		Address:     organization.Address.ToEntity(),
		Type:        organization.Type,
		CreatorID:   organization.CreatorID,
		CreatedAt:   organization.CreatedAt,
		UpdatedAt:   organization.UpdatedAt,
	}
}
