package models

import (
	"relif/bff/entities"
	"time"
)

type ProductType struct {
	ID             string    `bson:"_id"`
	Name           string    `bson:"name"`
	Description    string    `bson:"description"`
	Brand          string    `bson:"brand"`
	Category       string    `bson:"category"`
	OrganizationID string    `bson:"organization_id"`
	TotalInStock   int       `bson:"total_in_stock"`
	CreatedAt      time.Time `bson:"created_at"`
	UpdatedAt      time.Time `bson:"updated_at"`
}

func (productType *ProductType) ToEntity() entities.ProductType {
	return entities.ProductType{
		ID:             productType.ID,
		Name:           productType.Name,
		Description:    productType.Description,
		Brand:          productType.Brand,
		Category:       productType.Category,
		OrganizationID: productType.OrganizationID,
		TotalInStock:   productType.TotalInStock,
		CreatedAt:      productType.CreatedAt,
		UpdatedAt:      productType.UpdatedAt,
	}
}
