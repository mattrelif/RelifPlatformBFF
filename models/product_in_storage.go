package models

import "relif/platform-bff/entities"

type ProductInStorage struct {
	ID             string   `bson:"_id,omitempty"`
	OrganizationID string   `bson:"organization_id,omitempty"`
	Location       Location `bson:"location,omitempty"`
	ProductTypeID  string   `bson:"product_type_id,omitempty"`
}

func NewProductInStorage(entity entities.ProductInStorage) ProductInStorage {
	return ProductInStorage{
		ID:             entity.ID,
		OrganizationID: entity.OrganizationID,
		Location:       NewLocation(entity.Location),
		ProductTypeID:  entity.ProductTypeID,
	}
}

func NewUpdatedProductInStorage(entity entities.ProductInStorage) ProductInStorage {
	return ProductInStorage{
		Location: NewLocation(entity.Location),
	}
}
