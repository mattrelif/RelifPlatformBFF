package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"relif/bff/entities"
	"time"
)

type FindProductAllocation struct {
	ID             string       `bson:"_id,omitempty"`
	ProductTypeID  string       `bson:"product_type_id,omitempty"`
	ProductType    ProductType  `bson:"product_type,omitempty"`
	Type           string       `bson:"type,omitempty"`
	FromHousingID  string       `bson:"from_housing_id,omitempty"`
	FromHousing    Housing      `bson:"from_housing,omitempty"`
	ToHousingID    string       `bson:"to_housing_id,omitempty"`
	ToHousing      Housing      `bson:"to_housing,omitempty"`
	OrganizationID string       `bson:"organization_id,omitempty"`
	Organization   Organization `bson:"organization,omitempty"`
	Quantity       int          `bson:"quantity,omitempty"`
	AuditorID      string       `bson:"auditor_id,omitempty"`
	Auditor        User         `bson:"auditor,omitempty"`
	CreatedAt      time.Time    `bson:"created_at,omitempty"`
}

func (productAllocation *FindProductAllocation) ToEntity() entities.ProductAllocation {
	return entities.ProductAllocation{
		ID:             productAllocation.ID,
		ProductTypeID:  productAllocation.ProductTypeID,
		ProductType:    productAllocation.ProductType.ToEntity(),
		Type:           productAllocation.Type,
		FromHousingID:  productAllocation.FromHousingID,
		FromHousing:    productAllocation.FromHousing.ToEntity(),
		ToHousingID:    productAllocation.ToHousingID,
		ToHousing:      productAllocation.ToHousing.ToEntity(),
		OrganizationID: productAllocation.OrganizationID,
		Organization:   productAllocation.Organization.ToEntity(),
		Quantity:       productAllocation.Quantity,
		AuditorID:      productAllocation.AuditorID,
		Auditor:        productAllocation.Auditor.ToEntity(),
		CreatedAt:      productAllocation.CreatedAt,
	}
}

type ProductAllocation struct {
	ID             string    `bson:"_id,omitempty"`
	ProductTypeID  string    `bson:"product_type_id,omitempty"`
	Type           string    `bson:"type,omitempty"`
	FromHousingID  string    `bson:"from_housing_id,omitempty"`
	ToHousingID    string    `bson:"to_housing_id,omitempty"`
	OrganizationID string    `bson:"organization_id,omitempty"`
	Quantity       int       `bson:"quantity,omitempty"`
	AuditorID      string    `bson:"auditor_id,omitempty"`
	CreatedAt      time.Time `bson:"created_at,omitempty"`
}

func (productAllocation *ProductAllocation) ToEntity() entities.ProductAllocation {
	return entities.ProductAllocation{
		ID:             productAllocation.ID,
		ProductTypeID:  productAllocation.ProductTypeID,
		Type:           productAllocation.Type,
		FromHousingID:  productAllocation.FromHousingID,
		ToHousingID:    productAllocation.ToHousingID,
		OrganizationID: productAllocation.OrganizationID,
		Quantity:       productAllocation.Quantity,
		AuditorID:      productAllocation.AuditorID,
		CreatedAt:      productAllocation.CreatedAt,
	}
}

func NewProductAllocation(entity entities.ProductAllocation) ProductAllocation {
	return ProductAllocation{
		ID:             primitive.NewObjectID().Hex(),
		ProductTypeID:  entity.ProductTypeID,
		Type:           entity.Type,
		FromHousingID:  entity.FromHousingID,
		ToHousingID:    entity.ToHousingID,
		OrganizationID: entity.OrganizationID,
		Quantity:       entity.Quantity,
		AuditorID:      entity.AuditorID,
		CreatedAt:      time.Now(),
	}
}
