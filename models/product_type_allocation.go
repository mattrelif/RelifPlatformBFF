package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"relif/platform-bff/entities"
	"time"
)

type FindProductTypeAllocation struct {
	ID             string       `bson:"_id,omitempty"`
	ProductTypeID  string       `bson:"product_type_id,omitempty"`
	Type           string       `bson:"type,omitempty"`
	From           FindLocation `bson:"from,omitempty"`
	To             FindLocation `bson:"to,omitempty"`
	OrganizationID string       `bson:"organization_id,omitempty"`
	Quantity       int          `bson:"quantity,omitempty"`
	CreatedAt      time.Time    `bson:"created_at,omitempty"`
}

func (productTypeAllocation *FindProductTypeAllocation) ToEntity() entities.ProductTypeAllocation {
	return entities.ProductTypeAllocation{
		ID:             productTypeAllocation.ID,
		ProductTypeID:  productTypeAllocation.ProductTypeID,
		Type:           productTypeAllocation.Type,
		From:           productTypeAllocation.From.ToEntity(),
		To:             productTypeAllocation.To.ToEntity(),
		OrganizationID: productTypeAllocation.OrganizationID,
		Quantity:       productTypeAllocation.Quantity,
		CreatedAt:      productTypeAllocation.CreatedAt,
	}
}

type ProductTypeAllocation struct {
	ID             string    `bson:"_id,omitempty"`
	ProductTypeID  string    `bson:"product_type_id,omitempty"`
	Type           string    `bson:"type,omitempty"`
	From           Location  `bson:"from,omitempty"`
	To             Location  `bson:"to,omitempty"`
	OrganizationID string    `bson:"organization_id,omitempty"`
	Quantity       int       `bson:"quantity,omitempty"`
	CreatedAt      time.Time `bson:"created_at,omitempty"`
}

func (productTypeAllocation *ProductTypeAllocation) ToEntity() entities.ProductTypeAllocation {
	return entities.ProductTypeAllocation{
		ID:             productTypeAllocation.ID,
		ProductTypeID:  productTypeAllocation.ProductTypeID,
		Type:           productTypeAllocation.Type,
		From:           productTypeAllocation.From.ToEntity(),
		To:             productTypeAllocation.To.ToEntity(),
		OrganizationID: productTypeAllocation.OrganizationID,
		Quantity:       productTypeAllocation.Quantity,
		CreatedAt:      productTypeAllocation.CreatedAt,
	}
}

func NewProductTypeAllocation(entity entities.ProductTypeAllocation) ProductTypeAllocation {
	return ProductTypeAllocation{
		ID:             primitive.NewObjectID().Hex(),
		ProductTypeID:  entity.ProductTypeID,
		Type:           entity.Type,
		From:           NewLocation(entity.From),
		To:             NewLocation(entity.To),
		OrganizationID: entity.OrganizationID,
		Quantity:       entity.Quantity,
		CreatedAt:      time.Now(),
	}
}
