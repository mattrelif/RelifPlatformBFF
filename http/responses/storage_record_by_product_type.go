package responses

import "relif/platform-bff/entities"

type StorageRecordsByProductType []StorageRecordByProductType

type StorageRecordByProductType struct {
	ID       string   `json:"id,omitempty"`
	Location Location `json:"location,omitempty"`
	Quantity int      `json:"quantity,omitempty"`
}

func NewStorageRecordsByProductType(entityList []entities.StorageRecord) StorageRecordsByProductType {
	res := make(StorageRecordsByProductType, 0)

	for _, entity := range entityList {
		res = append(res, NewStorageRecordByProductType(entity))
	}

	return res
}

func NewStorageRecordByProductType(entity entities.StorageRecord) StorageRecordByProductType {
	return StorageRecordByProductType{
		ID:       entity.ID,
		Location: NewLocation(entity.Location),
		Quantity: entity.Quantity,
	}
}
