package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"relif/platform-bff/entities"
	"relif/platform-bff/utils"
	"time"
)

type Organization struct {
	ID               string    `bson:"_id,omitempty"`
	Name             string    `bson:"name,omitempty"`
	Description      string    `bson:"description,omitempty"`
	Logo             string    `bson:"logo,omitempty"`
	AreasOfWork      []string  `bson:"areas_of_work,omitempty"`
	Address          Address   `bson:"address,omitempty"`
	Type             string    `bson:"type,omitempty"`
	Status           string    `bson:"status,omitempty"`
	OwnerID          string    `bson:"owner_id,omitempty"`
	AccessGrantedIDs []string  `bson:"access_granted_ids,omitempty"`
	CreatedAt        time.Time `bson:"created_at,omitempty"`
	UpdatedAt        time.Time `bson:"updated_at,omitempty"`
}

func (organization *Organization) ToEntity() entities.Organization {
	return entities.Organization{
		ID:               organization.ID,
		Name:             organization.Name,
		Description:      organization.Description,
		Logo:             organization.Logo,
		AreasOfWork:      organization.AreasOfWork,
		Address:          organization.Address.ToEntity(),
		Type:             organization.Type,
		Status:           organization.Status,
		OwnerID:          organization.OwnerID,
		AccessGrantedIDs: organization.AccessGrantedIDs,
		CreatedAt:        organization.CreatedAt,
		UpdatedAt:        organization.UpdatedAt,
	}
}

func NewOrganization(entity entities.Organization) Organization {
	return Organization{
		ID:          primitive.NewObjectID().Hex(),
		Name:        entity.Name,
		Description: entity.Description,
		Logo:        entity.Logo,
		AreasOfWork: entity.AreasOfWork,
		Address:     NewAddress(entity.Address),
		Type:        utils.ManagerOrganizationType,
		Status:      utils.ActiveStatus,
		OwnerID:     entity.OwnerID,
		CreatedAt:   time.Now(),
	}
}

func NewUpdatedOrganization(entity entities.Organization) Organization {
	return Organization{
		Name:             entity.Name,
		Description:      entity.Description,
		Logo:             entity.Logo,
		AreasOfWork:      entity.AreasOfWork,
		Address:          NewAddress(entity.Address),
		Type:             entity.Type,
		OwnerID:          entity.OwnerID,
		Status:           entity.Status,
		AccessGrantedIDs: entity.AccessGrantedIDs,
		UpdatedAt:        time.Now(),
	}
}
