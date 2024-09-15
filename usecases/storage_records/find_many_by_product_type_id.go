package storage_records

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type FindManyByProductTypeID interface {
	Execute(actor entities.User, productTypeID string) ([]entities.StorageRecord, error)
}

type findManyByProductTypeIDImpl struct {
	productTypesRepository   repositories.ProductTypes
	organizationsRepository  repositories.Organizations
	storageRecordsRepository repositories.StorageRecords
}

func NewFindManyByProductTypeID(
	productTypesRepository repositories.ProductTypes,
	organizationsRepository repositories.Organizations,
	storageRecordsRepository repositories.StorageRecords,
) FindManyByProductTypeID {
	return &findManyByProductTypeIDImpl{
		productTypesRepository:   productTypesRepository,
		organizationsRepository:  organizationsRepository,
		storageRecordsRepository: storageRecordsRepository,
	}
}

func (uc *findManyByProductTypeIDImpl) Execute(actor entities.User, productTypeID string) ([]entities.StorageRecord, error) {
	productType, err := uc.productTypesRepository.FindOneByID(productTypeID)

	if err != nil {
		return nil, err
	}

	organization, err := uc.organizationsRepository.FindOneByID(productType.OrganizationID)

	if err != nil {
		return nil, err
	}

	if err = guards.IsOrganizationAdmin(actor, organization); err != nil {
		return nil, err
	}

	return uc.storageRecordsRepository.FindManyByProductTypeID(productTypeID)
}
