package product_type_alloctions

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type FindManyByProductTypeIDPaginated interface {
	Execute(actor entities.User, productTypeID string, offset, limit int64) (int64, []entities.ProductTypeAllocation, error)
}

type findManyByProductTypeIDPaginatedImpl struct {
	productTypesRepository           repositories.ProductTypes
	productTypeAllocationsRepository repositories.ProductTypeAllocations
	organizationsRepository          repositories.Organizations
}

func NewFindManyByProductTypeIDPaginated(
	productTypesRepository repositories.ProductTypes,
	productTypeAllocationsRepository repositories.ProductTypeAllocations,
	organizationsRepository repositories.Organizations,
) FindManyByProductTypeIDPaginated {
	return &findManyByProductTypeIDPaginatedImpl{
		productTypesRepository:           productTypesRepository,
		productTypeAllocationsRepository: productTypeAllocationsRepository,
		organizationsRepository:          organizationsRepository,
	}
}

func (uc *findManyByProductTypeIDPaginatedImpl) Execute(actor entities.User, productTypeID string, offset, limit int64) (int64, []entities.ProductTypeAllocation, error) {
	productType, err := uc.productTypesRepository.FindOneByID(productTypeID)

	if err != nil {
		return 0, nil, err
	}

	organization, err := uc.organizationsRepository.FindOneByID(productType.OrganizationID)

	if err != nil {
		return 0, nil, err
	}

	if err = guards.IsOrganizationAdmin(actor, organization); err != nil {
		return 0, nil, err
	}

	return uc.productTypeAllocationsRepository.FindManyByProductTypeIDPaginated(productTypeID, offset, limit)
}
