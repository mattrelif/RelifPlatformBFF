package product_types

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type DeleteOneByID interface {
	Execute(actor entities.User, id string) error
}

type deleteOneByIDImpl struct {
	productTypesRepository  repositories.ProductTypes
	organizationsRepository repositories.Organizations
}

func NewDeleteOneByID(
	productTypesRepository repositories.ProductTypes,
	organizationsRepository repositories.Organizations,
) DeleteOneByID {
	return &deleteOneByIDImpl{
		productTypesRepository:  productTypesRepository,
		organizationsRepository: organizationsRepository,
	}
}

func (uc *deleteOneByIDImpl) Execute(actor entities.User, id string) error {
	productType, err := uc.productTypesRepository.FindOneByID(id)

	if err != nil {
		return err
	}

	organization, err := uc.organizationsRepository.FindOneByID(productType.OrganizationID)

	if err != nil {
		return err
	}

	if err = guards.IsOrganizationAdmin(actor, organization); err != nil {
		return err
	}

	return uc.productTypesRepository.DeleteOneByID(productType.ID)
}
