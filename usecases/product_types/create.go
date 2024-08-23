package product_types

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type Create interface {
	Execute(actor entities.User, organizationID string, data entities.ProductType) (entities.ProductType, error)
}

type createImpl struct {
	organizationsRepository repositories.Organizations
	productTypesRepository  repositories.ProductTypes
}

func NewCreate(
	organizationsRepository repositories.Organizations,
	productTypesRepository repositories.ProductTypes,
) Create {
	return &createImpl{
		organizationsRepository: organizationsRepository,
		productTypesRepository:  productTypesRepository,
	}
}

func (uc *createImpl) Execute(actor entities.User, organizationID string, data entities.ProductType) (entities.ProductType, error) {
	organization, err := uc.organizationsRepository.FindOneByID(organizationID)

	if err != nil {
		return entities.ProductType{}, err
	}

	if err = guards.IsOrganizationAdmin(actor, organization); err != nil {
		return entities.ProductType{}, err
	}

	data.OrganizationID = organization.ID

	return uc.productTypesRepository.Create(data)
}
