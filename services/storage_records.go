package services

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/repositories"
)

type StorageRecords interface {
	Create(data entities.StorageRecord) error
	FindOneByProductTypeIDAndLocation(productTypeID string, location entities.Location) (entities.StorageRecord, error)
	UpdateOneByID(id string, data entities.StorageRecord) error
}

type storageRecordsImpl struct {
	repository repositories.StorageRecords
}

func NewStorageRecords(repository repositories.StorageRecords) StorageRecords {
	return &storageRecordsImpl{
		repository: repository,
	}
}

func (service *storageRecordsImpl) Create(data entities.StorageRecord) error {
	return service.repository.Create(data)
}

func (service *storageRecordsImpl) FindOneByProductTypeIDAndLocation(productTypeID string, location entities.Location) (entities.StorageRecord, error) {
	return service.repository.FindOneByProductTypeIDAndLocation(productTypeID, location)
}

func (service *storageRecordsImpl) UpdateOneByID(id string, data entities.StorageRecord) error {
	return service.repository.UpdateOneByID(id, data)
}
