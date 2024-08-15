package responses

import "relif/platform-bff/entities"

type StorageRecordsByLocation []StorageRecordByLocation

type StorageRecordByLocation struct {
	ID            string      `json:"id,omitempty"`
	ProductTypeID string      `json:"product_type_id,omitempty"`
	ProductType   ProductType `json:"product_type,omitempty"`
	Quantity      int         `json:"quantity,omitempty"`
}

func NewStorageRecordsByLocation(entityList []entities.StorageRecord) StorageRecordsByLocation {
	res := make(StorageRecordsByLocation, 0)

	for _, entity := range entityList {
		res = append(res, NewStorageRecordByLocation(entity))
	}

	return res
}

func NewStorageRecordByLocation(entity entities.StorageRecord) StorageRecordByLocation {
	return StorageRecordByLocation{
		ID:            entity.ID,
		ProductTypeID: entity.ProductTypeID,
		ProductType:   NewProductType(entity.ProductType),
		Quantity:      entity.Quantity,
	}
}
