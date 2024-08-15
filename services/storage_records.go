package services

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/repositories"
	"relif/platform-bff/utils"
)

type StorageRecords interface {
	Create(data entities.StorageRecord) error
	FindOneByProductTypeIDAndLocation(productTypeID string, location entities.Location) (entities.StorageRecord, error)
	FindManyByHousingID(housingID string, offset, limit int64) (int64, []entities.StorageRecord, error)
	FindManyByOrganizationID(organizationID string, offset, limit int64) (int64, []entities.StorageRecord, error)
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

func (service *storageRecordsImpl) FindManyByHousingID(housingID string, offset, limit int64) (int64, []entities.StorageRecord, error) {
	location := entities.Location{
		ID:   housingID,
		Type: utils.HousingLocationType,
	}

	return service.repository.FindManyByLocation(location, offset, limit)
}

func (service *storageRecordsImpl) FindManyByOrganizationID(organizationID string, offset, limit int64) (int64, []entities.StorageRecord, error) {
	location := entities.Location{
		ID:   organizationID,
		Type: utils.HousingLocationType,
	}

	return service.repository.FindManyByLocation(location, offset, limit)
}
