package services

import (
	"relif/bff/entities"
	"relif/bff/repositories"
)

type ProductTypes interface {
	Create(user entities.User, data entities.ProductType) (entities.ProductType, error)
	FindManyByOrganizationId(organizationId string, limit, offset int64) (int64, []entities.ProductType, error)
	FindOneById(id string) (entities.ProductType, error)
	UpdateOneById(id string, data entities.ProductType) error
	IncreaseTotalInStock(id string, amount int) error
	DeleteOneById(id string) error
}

type productTypesImpl struct {
	repository repositories.ProductTypes
}

func NewProductTypes(repository repositories.ProductTypes) ProductTypes {
	return &productTypesImpl{
		repository: repository,
	}
}

func (service *productTypesImpl) Create(user entities.User, data entities.ProductType) (entities.ProductType, error) {
	data.OrganizationID = user.OrganizationID
	return service.repository.Create(data)
}

func (service *productTypesImpl) FindManyByOrganizationId(organizationId string, limit, offset int64) (int64, []entities.ProductType, error) {
	return service.repository.FindManyByOrganizationId(organizationId, limit, offset)
}

func (service *productTypesImpl) FindOneById(id string) (entities.ProductType, error) {
	return service.repository.FindOneById(id)
}

func (service *productTypesImpl) UpdateOneById(id string, data entities.ProductType) error {
	return service.repository.UpdateOneById(id, data)
}

func (service *productTypesImpl) IncreaseTotalInStock(id string, amount int) error {
	return service.repository.IncreaseTotalInStock(id, amount)
}

func (service *productTypesImpl) DeleteOneById(id string) error {
	return service.repository.DeleteOneById(id)
}
