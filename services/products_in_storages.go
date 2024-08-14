package services

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/repositories"
)

type ProductsInStorages interface {
	CreateMany(data entities.ProductInStorage, quantity int) error
	FindManyIDsByLocation(location entities.Location, quantity int) ([]interface{}, error)
	UpdateManyByIDs(ids []interface{}, data entities.ProductInStorage) error
	DeleteManyByIDs(ids []interface{}) error
}

type productsInStoragesImpl struct {
	repository repositories.ProductsInStorages
}

func NewProductsInStorages(repository repositories.ProductsInStorages) ProductsInStorages {
	return &productsInStoragesImpl{
		repository: repository,
	}
}

func (service *productsInStoragesImpl) CreateMany(data entities.ProductInStorage, quantity int) error {
	return service.repository.CreateMany(data, quantity)
}

func (service *productsInStoragesImpl) FindManyIDsByLocation(location entities.Location, quantity int) ([]interface{}, error) {
	return service.repository.FindManyIDsByLocation(location, quantity)
}

func (service *productsInStoragesImpl) UpdateManyByIDs(ids []interface{}, data entities.ProductInStorage) error {
	return service.repository.UpdateManyByIDs(ids, data)
}

func (service *productsInStoragesImpl) DeleteManyByIDs(ids []interface{}) error {
	return service.repository.DeleteManyByIDs(ids)
}
