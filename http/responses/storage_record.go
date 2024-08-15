package responses

import "relif/platform-bff/entities"

type StorageRecords []StorageRecord

type StorageRecord struct {
	Location Location `json:"location,omitempty"`
	Quantity int      `json:"quantity,omitempty"`
}

func NewStorageRecords(entityList []entities.StorageRecord) StorageRecords {
	res := make(StorageRecords, 0)

	for _, entity := range entityList {
		res = append(res, NewStorageRecord(entity))
	}

	return res
}

func NewStorageRecord(entity entities.StorageRecord) StorageRecord {
	return StorageRecord{
		Location: NewLocation(entity.Location),
		Quantity: entity.Quantity,
	}
}
