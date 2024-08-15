package models

import "relif/platform-bff/entities"

type FindStorageRecord struct {
	Location FindLocation `bson:"location,omitempty"`
	Quantity int          `bson:"quantity,omitempty"`
}

func (record *FindStorageRecord) ToEntity() entities.StorageRecord {
	return entities.StorageRecord{
		Location: record.Location.ToEntity(),
		Quantity: record.Quantity,
	}
}

type StorageRecord struct {
	ID            string   `bson:"_id,omitempty"`
	Location      Location `bson:"location,omitempty"`
	ProductTypeID string   `bson:"product_type_id,omitempty"`
	Quantity      int      `bson:"quantity,omitempty"`
}

func (record *StorageRecord) ToEntity() entities.StorageRecord {
	return entities.StorageRecord{
		ID:            record.ID,
		Location:      record.Location.ToEntity(),
		ProductTypeID: record.ProductTypeID,
		Quantity:      record.Quantity,
	}
}

func NewStorageRecord(entity entities.StorageRecord) StorageRecord {
	return StorageRecord{
		ID:            entity.ID,
		Location:      NewLocation(entity.Location),
		ProductTypeID: entity.ProductTypeID,
		Quantity:      entity.Quantity,
	}
}

func NewUpdatedStorageRecord(entity entities.StorageRecord) StorageRecord {
	return StorageRecord{
		Quantity: entity.Quantity,
	}
}
