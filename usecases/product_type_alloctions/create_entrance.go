package product_type_alloctions

import (
	"fmt"
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
	"relif/platform-bff/utils"
)

type CreateEntrance interface {
	Execute(actor entities.User, productTypeID string, data entities.ProductTypeAllocation) (entities.ProductTypeAllocation, error)
}

type createEntranceImpl struct {
	productTypeAllocationsRepository repositories.ProductTypeAllocations
	productTypesRepository           repositories.ProductTypes
	organizationsRepository          repositories.Organizations
	storageRecordRepository          repositories.StorageRecords
}

func NewCreateEntrance(
	productTypeAllocationsRepository repositories.ProductTypeAllocations,
	productTypesRepository repositories.ProductTypes,
	organizationsRepository repositories.Organizations,
	storageRecordRepository repositories.StorageRecords,
) CreateEntrance {
	return &createEntranceImpl{
		productTypeAllocationsRepository: productTypeAllocationsRepository,
		productTypesRepository:           productTypesRepository,
		organizationsRepository:          organizationsRepository,
		storageRecordRepository:          storageRecordRepository,
	}
}

func (uc *createEntranceImpl) Execute(actor entities.User, productTypeID string, data entities.ProductTypeAllocation) (entities.ProductTypeAllocation, error) {
	productType, err := uc.productTypesRepository.FindOneByID(productTypeID)

	if err != nil {
		return entities.ProductTypeAllocation{}, err
	}

	organization, err := uc.organizationsRepository.FindOneByID(productType.OrganizationID)

	if err != nil {
		return entities.ProductTypeAllocation{}, err
	}

	if err = guards.IsOrganizationAdmin(actor, organization); err != nil {
		return entities.ProductTypeAllocation{}, err
	}

	record, err := uc.storageRecordRepository.FindOneByProductTypeIDAndLocation(productType.ID, data.To)

	if err != nil {
		return entities.ProductTypeAllocation{}, err
	}

	if record.ID != "" {
		record.Quantity = record.Quantity + data.Quantity

		fmt.Println(record.Quantity)

		if err = uc.storageRecordRepository.UpdateOneByID(record.ID, record); err != nil {
			return entities.ProductTypeAllocation{}, err
		}
	} else {
		record = entities.StorageRecord{
			Quantity:      data.Quantity,
			Location:      data.To,
			ProductTypeID: productType.ID,
		}

		if err = uc.storageRecordRepository.Create(record); err != nil {
			return entities.ProductTypeAllocation{}, err
		}
	}

	data.Type = utils.EntranceType
	data.OrganizationID = actor.Organization.ID

	return uc.productTypeAllocationsRepository.Create(data)
}
