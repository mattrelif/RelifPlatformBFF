package product_type_alloctions

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
	"relif/platform-bff/utils"
)

type CreateReallocation interface {
	Execute(actor entities.User, productTypeID string, data entities.ProductTypeAllocation) (entities.ProductTypeAllocation, error)
}

type createReallocationImpl struct {
	productTypeAllocationsRepository repositories.ProductTypeAllocations
	productTypesRepository           repositories.ProductTypes
	organizationsRepository          repositories.Organizations
	storageRecordRepository          repositories.StorageRecords
}

func NewCreateReallocation(
	productTypeAllocationsRepository repositories.ProductTypeAllocations,
	productTypesRepository repositories.ProductTypes,
	organizationsRepository repositories.Organizations,
	storageRecordRepository repositories.StorageRecords,
) CreateReallocation {
	return &createReallocationImpl{
		productTypeAllocationsRepository: productTypeAllocationsRepository,
		productTypesRepository:           productTypesRepository,
		organizationsRepository:          organizationsRepository,
		storageRecordRepository:          storageRecordRepository,
	}
}

func (uc *createReallocationImpl) Execute(actor entities.User, productTypeID string, data entities.ProductTypeAllocation) (entities.ProductTypeAllocation, error) {
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

	fromRecord, err := uc.storageRecordRepository.FindOneByProductTypeIDAndLocation(productType.ID, data.From)

	if err != nil {
		return entities.ProductTypeAllocation{}, err
	}

	if fromRecord.ID != "" {
		if err = uc.storageRecordRepository.DecreaseQuantityOfOneByID(fromRecord.ID, data.Quantity); err != nil {
			return entities.ProductTypeAllocation{}, err
		}
	} else {
		return entities.ProductTypeAllocation{}, utils.ErrStorageRecordNotFound
	}

	toRecord, err := uc.storageRecordRepository.FindOneByProductTypeIDAndLocation(productType.ID, data.To)

	if err != nil {
		return entities.ProductTypeAllocation{}, err
	}

	if toRecord.ID != "" {
		if err = uc.storageRecordRepository.IncreaseQuantityOfOneByID(toRecord.ID, data.Quantity); err != nil {
			return entities.ProductTypeAllocation{}, err
		}
	} else {
		toRecord = entities.StorageRecord{
			Quantity:      data.Quantity,
			Location:      data.To,
			ProductTypeID: productType.ID,
		}

		if err = uc.storageRecordRepository.Create(toRecord); err != nil {
			return entities.ProductTypeAllocation{}, err
		}
	}

	data.Type = utils.ReallocationType
	data.OrganizationID = actor.Organization.ID

	return uc.productTypeAllocationsRepository.Create(data)
}
