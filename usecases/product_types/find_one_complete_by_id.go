package product_types

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type FindOneCompleteByID interface {
	Execute(actor entities.User, id string) (entities.ProductType, error)
}

type findOneCompleteByIDImpl struct {
	productTypesRepository repositories.ProductTypes
}

func NewFindOneCompleteByID(productTypesRepository repositories.ProductTypes) FindOneCompleteByID {
	return &findOneCompleteByIDImpl{
		productTypesRepository: productTypesRepository,
	}
}

func (uc *findOneCompleteByIDImpl) Execute(actor entities.User, id string) (entities.ProductType, error) {
	productType, err := uc.productTypesRepository.FindOneCompleteByID(id)

	if err != nil {
		return entities.ProductType{}, err
	}

	if err = guards.HasAccessToOrganizationData(actor, productType.Organization); err != nil {
		return entities.ProductType{}, err
	}

	return productType, nil
}
