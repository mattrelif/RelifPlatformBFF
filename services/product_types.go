package services

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/repositories"
)

type ProductTypes interface {
	Create(user entities.User, data entities.ProductType) (entities.ProductType, error)
	FindManyByOrganizationID(organizationID string, limit, offset int64) (int64, []entities.ProductType, error)
	FindOneByID(id string) (entities.ProductType, error)
	UpdateOneByID(id string, data entities.ProductType) error
	IncreaseTotalInStock(id string, amount int) error
	DeleteOneByID(id string) error
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

func (service *productTypesImpl) FindManyByOrganizationID(organizationID string, limit, offset int64) (int64, []entities.ProductType, error) {
	return service.repository.FindManyByOrganizationID(organizationID, limit, offset)
}

func (service *productTypesImpl) FindOneByID(id string) (entities.ProductType, error) {
	return service.repository.FindOneByID(id)
}

func (service *productTypesImpl) UpdateOneByID(id string, data entities.ProductType) error {
	return service.repository.UpdateOneByID(id, data)
}

func (service *productTypesImpl) IncreaseTotalInStock(id string, amount int) error {
	return service.repository.IncreaseTotalInStock(id, amount)
}

func (service *productTypesImpl) DeleteOneByID(id string) error {
	return service.repository.DeleteOneByID(id)
}
