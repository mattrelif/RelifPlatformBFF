package product_types

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type UpdateOneByID interface {
	Execute(actor entities.User, id string, data entities.ProductType) error
}

type updateOneByIDImpl struct {
	productTypesRepository  repositories.ProductTypes
	organizationsRepository repositories.Organizations
}

func NewUpdateOneByID(
	productTypesRepository repositories.ProductTypes,
	organizationsRepository repositories.Organizations,
) UpdateOneByID {
	return &updateOneByIDImpl{
		productTypesRepository:  productTypesRepository,
		organizationsRepository: organizationsRepository,
	}
}

func (uc *updateOneByIDImpl) Execute(actor entities.User, id string, data entities.ProductType) error {
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

	return uc.productTypesRepository.UpdateOneByID(productType.ID, data)
}
