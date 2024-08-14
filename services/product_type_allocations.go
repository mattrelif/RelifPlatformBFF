package services

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/repositories"
	"relif/platform-bff/utils"
)

type ProductTypeAllocations interface {
	CreateEntrance(productTypeID string, data entities.ProductTypeAllocation) (entities.ProductTypeAllocation, error)
	CreateReallocation(productTypeID string, data entities.ProductTypeAllocation) (entities.ProductTypeAllocation, error)
}

type productTypeAllocationsImpl struct {
	repository                repositories.ProductTypeAllocations
	productTypesService       ProductTypes
	productsInStoragesService ProductsInStorages
}

func NewProductTypeAllocations(
	repository repositories.ProductTypeAllocations,
	productTypesService ProductTypes,
	productsInStoragesService ProductsInStorages,
) ProductTypeAllocations {
	return &productTypeAllocationsImpl{
		repository:                repository,
		productTypesService:       productTypesService,
		productsInStoragesService: productsInStoragesService,
	}
}

func (service *productTypeAllocationsImpl) CreateEntrance(productTypeID string, data entities.ProductTypeAllocation) (entities.ProductTypeAllocation, error) {
	productType, err := service.productTypesService.FindOneByID(productTypeID)

	if err != nil {
		return entities.ProductTypeAllocation{}, err
	}

	data.OrganizationID = productType.OrganizationID
	data.Type = utils.EntranceType

	product := entities.ProductInStorage{
		ProductTypeID:  data.ProductTypeID,
		Location:       data.To,
		OrganizationID: data.OrganizationID,
	}

	if err = service.productsInStoragesService.CreateMany(product, data.Quantity); err != nil {
		return entities.ProductTypeAllocation{}, err
	}

	return service.repository.Create(data)
}

func (service *productTypeAllocationsImpl) CreateReallocation(productTypeID string, data entities.ProductTypeAllocation) (entities.ProductTypeAllocation, error) {
	productType, err := service.productTypesService.FindOneByID(productTypeID)

	if err != nil {
		return entities.ProductTypeAllocation{}, err
	}

	data.OrganizationID = productType.OrganizationID
	data.Type = utils.ReallocationType

	ids, err := service.productsInStoragesService.FindManyIDsByLocation(data.From, data.Quantity)

	if err != nil {
		return entities.ProductTypeAllocation{}, err
	}

	product := entities.ProductInStorage{
		Location: data.To,
	}

	if err = service.productsInStoragesService.UpdateManyByIDs(ids, product); err != nil {
		return entities.ProductTypeAllocation{}, err
	}

	return service.repository.Create(data)
}
