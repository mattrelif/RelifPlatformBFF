package donations

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type FindManyByProductTypeIDPaginated interface {
	Execute(actor entities.User, productTypeID string, offset, limit int64) (int64, []entities.Donation, error)
}

type findManyByProductTypeIDPaginatedImpl struct {
	donationsRepository     repositories.Donations
	productTypesRepository  repositories.ProductTypes
	organizationsRepository repositories.Organizations
}

func NewFindManyByProductTypeIDPaginated(
	donationsRepository repositories.Donations,
	productTypesRepository repositories.ProductTypes,
	organizationsRepository repositories.Organizations,
) FindManyByProductTypeIDPaginated {
	return &findManyByProductTypeIDPaginatedImpl{
		donationsRepository:     donationsRepository,
		productTypesRepository:  productTypesRepository,
		organizationsRepository: organizationsRepository,
	}
}

func (uc *findManyByProductTypeIDPaginatedImpl) Execute(actor entities.User, productTypeID string, offset, limit int64) (int64, []entities.Donation, error) {
	productType, err := uc.productTypesRepository.FindOneByID(productTypeID)

	if err != nil {
		return 0, nil, err
	}

	organization, err := uc.organizationsRepository.FindOneByID(productType.OrganizationID)

	if err != nil {
		return 0, nil, err
	}

	if err = guards.HasAccessToOrganizationData(actor, organization); err != nil {
		return 0, nil, err
	}

	return uc.donationsRepository.FindManyByProductTypeIDPaginated(productTypeID, offset, limit)
}
