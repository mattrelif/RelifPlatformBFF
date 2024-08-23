package product_types

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type FindManyByOrganizationIDPaginated interface {
	Execute(actor entities.User, organizationID string, offset, limit int64) (int64, []entities.ProductType, error)
}

type findManyByOrganizationIDPaginatedImpl struct {
	productTypesRepository  repositories.ProductTypes
	organizationsRepository repositories.Organizations
}

func NewFindManyByOrganizationIDPaginated(
	productTypesRepository repositories.ProductTypes,
	organizationsRepository repositories.Organizations,
) FindManyByOrganizationIDPaginated {
	return &findManyByOrganizationIDPaginatedImpl{
		productTypesRepository:  productTypesRepository,
		organizationsRepository: organizationsRepository,
	}
}

func (uc *findManyByOrganizationIDPaginatedImpl) Execute(actor entities.User, organizationID string, offset, limit int64) (int64, []entities.ProductType, error) {
	organization, err := uc.organizationsRepository.FindOneByID(organizationID)

	if err != nil {
		return 0, nil, err
	}

	if err = guards.HasAccessToOrganizationData(actor, organization); err != nil {
		return 0, nil, err
	}

	return uc.productTypesRepository.FindManyByOrganizationIDPaginated(organizationID, offset, limit)
}
