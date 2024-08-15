package services

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/repositories"
	"relif/platform-bff/utils"
)

type ProductTypeAllocations interface {
	CreateEntrance(productTypeID string, data entities.ProductTypeAllocation) (entities.ProductTypeAllocation, error)
	CreateReallocation(productTypeID string, data entities.ProductTypeAllocation) (entities.ProductTypeAllocation, error)
}

type productTypeAllocationsImpl struct {
	repository            repositories.ProductTypeAllocations
	productTypesService   ProductTypes
	storageRecordsService StorageRecords
}

func NewProductTypeAllocations(
	repository repositories.ProductTypeAllocations,
	productTypesService ProductTypes,
	storageRecordsService StorageRecords,
) ProductTypeAllocations {
	return &productTypeAllocationsImpl{
		repository:            repository,
		productTypesService:   productTypesService,
		storageRecordsService: storageRecordsService,
	}
}

func (service *productTypeAllocationsImpl) CreateEntrance(productTypeID string, data entities.ProductTypeAllocation) (entities.ProductTypeAllocation, error) {
	productType, err := service.productTypesService.FindOneByID(productTypeID)

	if err != nil {
		return entities.ProductTypeAllocation{}, err
	}

	record, err := service.storageRecordsService.FindOneByProductTypeIDAndLocation(productType.ID, data.To)

	if err != nil {
		return entities.ProductTypeAllocation{}, err
	}

	if record.ID != "" {
		record.Quantity += data.Quantity

		if err = service.storageRecordsService.UpdateOneByID(record.ID, record); err != nil {
			return entities.ProductTypeAllocation{}, err
		}
	} else {
		record = entities.StorageRecord{
			ProductTypeID: productType.ID,
			Quantity:      data.Quantity,
			Location:      data.To,
		}

		if err = service.storageRecordsService.Create(record); err != nil {
			return entities.ProductTypeAllocation{}, err
		}
	}

	data.OrganizationID = productType.OrganizationID
	data.Type = utils.EntranceType

	return service.repository.Create(data)
}

func (service *productTypeAllocationsImpl) CreateReallocation(productTypeID string, data entities.ProductTypeAllocation) (entities.ProductTypeAllocation, error) {
	productType, err := service.productTypesService.FindOneByID(productTypeID)

	if err != nil {
		return entities.ProductTypeAllocation{}, err
	}

	fromRecord, err := service.storageRecordsService.FindOneByProductTypeIDAndLocation(productType.ID, data.From)

	if err != nil {
		return entities.ProductTypeAllocation{}, err
	}

	if fromRecord.ID != "" {
		fromRecord.Quantity -= data.Quantity

		if err = service.storageRecordsService.UpdateOneByID(fromRecord.ID, fromRecord); err != nil {
			return entities.ProductTypeAllocation{}, err
		}
	} else {
		return entities.ProductTypeAllocation{}, utils.ErrStorageRecordNotFound
	}

	toRecord, err := service.storageRecordsService.FindOneByProductTypeIDAndLocation(productType.ID, data.To)

	if err != nil {
		return entities.ProductTypeAllocation{}, err
	}

	if toRecord.ID != "" {
		toRecord.Quantity += data.Quantity

		if err = service.storageRecordsService.UpdateOneByID(toRecord.ID, toRecord); err != nil {
			return entities.ProductTypeAllocation{}, err
		}
	} else {
		toRecord = entities.StorageRecord{
			ProductTypeID: productType.ID,
			Quantity:      data.Quantity,
			Location:      data.To,
		}

		if err = service.storageRecordsService.Create(toRecord); err != nil {
			return entities.ProductTypeAllocation{}, err
		}
	}

	data.OrganizationID = productType.OrganizationID
	data.Type = utils.ReallocationType

	return service.repository.Create(data)
}
