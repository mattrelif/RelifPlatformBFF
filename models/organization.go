package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"relif/bff/entities"
	"relif/bff/utils"
	"time"
)

type Organization struct {
	ID          string    `bson:"_id,omitempty"`
	Name        string    `bson:"name,omitempty"`
	Description string    `bson:"description,omitempty"`
	Address     Address   `bson:"address,omitempty"`
	Type        string    `bson:"type,omitempty"`
	OwnerID     string    `bson:"owner_id,omitempty"`
	CreatedAt   time.Time `bson:"created_at,omitempty"`
	UpdatedAt   time.Time `bson:"updated_at,omitempty"`
}

func (organization *Organization) ToEntity() entities.Organization {
	return entities.Organization{
		ID:          organization.ID,
		Name:        organization.Name,
		Description: organization.Description,
		Address:     organization.Address.ToEntity(),
		Type:        organization.Type,
		OwnerID:     organization.OwnerID,
		CreatedAt:   organization.CreatedAt,
		UpdatedAt:   organization.UpdatedAt,
	}
}

func NewOrganization(entity entities.Organization) Organization {
	return Organization{
		ID:          primitive.NewObjectID().Hex(),
		Name:        entity.Name,
		Description: entity.Description,
		Address:     NewAddress(entity.Address),
		Type:        utils.ManagerOrganizationType,
		OwnerID:     entity.OwnerID,
		CreatedAt:   time.Now(),
	}
}

func NewUpdatedOrganization(entity entities.Organization) Organization {
	return Organization{
		Name:        entity.Name,
		Description: entity.Description,
		Address:     NewAddress(entity.Address),
		Type:        entity.Type,
		OwnerID:     entity.OwnerID,
		UpdatedAt:   time.Now(),
	}
}
