package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"relif/platform-bff/entities"
	"time"
)

type FindProductType struct {
	ID             string       `bson:"_id,omitempty"`
	Name           string       `bson:"name,omitempty"`
	Description    string       `bson:"description,omitempty"`
	Brand          string       `bson:"brand,omitempty"`
	Category       string       `bson:"category,omitempty"`
	OrganizationID string       `bson:"organization_id,omitempty"`
	Organization   Organization `bson:"organization,omitempty"`
	UnitType       string       `bson:"unit_type,omitempty"`
	TotalInStorage int          `bson:"total_in_storage,omitempty"`
	CreatedAt      time.Time    `bson:"created_at,omitempty"`
	UpdatedAt      time.Time    `bson:"updated_at,omitempty"`
}

func (productType *FindProductType) ToEntity() entities.ProductType {
	return entities.ProductType{
		ID:             productType.ID,
		Name:           productType.Name,
		Description:    productType.Description,
		Brand:          productType.Brand,
		Category:       productType.Category,
		OrganizationID: productType.OrganizationID,
		Organization:   productType.Organization.ToEntity(),
		UnitType:       productType.UnitType,
		TotalInStorage: productType.TotalInStorage,
		CreatedAt:      productType.CreatedAt,
		UpdatedAt:      productType.UpdatedAt,
	}
}

type ProductType struct {
	ID             string    `bson:"_id,omitempty"`
	Name           string    `bson:"name,omitempty"`
	Description    string    `bson:"description,omitempty"`
	Brand          string    `bson:"brand,omitempty"`
	Category       string    `bson:"category,omitempty"`
	OrganizationID string    `bson:"organization_id,omitempty"`
	UnitType       string    `bson:"unit_type,omitempty"`
	CreatedAt      time.Time `bson:"created_at,omitempty"`
	UpdatedAt      time.Time `bson:"updated_at,omitempty"`
}

func (productType *ProductType) ToEntity() entities.ProductType {
	return entities.ProductType{
		ID:             productType.ID,
		Name:           productType.Name,
		Description:    productType.Description,
		Brand:          productType.Brand,
		Category:       productType.Category,
		OrganizationID: productType.OrganizationID,
		UnitType:       productType.UnitType,
		CreatedAt:      productType.CreatedAt,
		UpdatedAt:      productType.UpdatedAt,
	}
}

func NewProductType(entity entities.ProductType) ProductType {
	return ProductType{
		ID:             primitive.NewObjectID().Hex(),
		Name:           entity.Name,
		Description:    entity.Description,
		Brand:          entity.Brand,
		Category:       entity.Category,
		OrganizationID: entity.OrganizationID,
		UnitType:       entity.UnitType,
		CreatedAt:      time.Now(),
	}
}

func NewUpdatedProductType(entity entities.ProductType) ProductType {
	return ProductType{
		Name:           entity.Name,
		Description:    entity.Description,
		Brand:          entity.Brand,
		Category:       entity.Category,
		OrganizationID: entity.OrganizationID,
		UnitType:       entity.UnitType,
		UpdatedAt:      time.Now(),
	}
}
