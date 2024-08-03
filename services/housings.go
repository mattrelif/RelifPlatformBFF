package services

import (
	"relif/bff/entities"
	"relif/bff/repositories"
	"relif/bff/utils"
)

type Housings interface {
	Create(user entities.User, housing entities.Housing) (entities.Housing, error)
	FindManyByOrganizationID(organizationId string, limit, offset int64) (int64, []entities.Housing, error)
	FindOneByID(id string) (entities.Housing, error)
	UpdateOneByID(id string, housing entities.Housing) error
	InactivateOneByID(id string) error
}

type housingsImpl struct {
	repository repositories.Housings
}

func NewHousings(repository repositories.Housings) Housings {
	return &housingsImpl{
		repository: repository,
	}
}

func (service *housingsImpl) Create(user entities.User, housing entities.Housing) (entities.Housing, error) {
	housing.OrganizationID = user.OrganizationID
	return service.repository.Create(housing)
}

func (service *housingsImpl) FindManyByOrganizationID(organizationId string, limit, offset int64) (int64, []entities.Housing, error) {
	return service.repository.FindManyByOrganizationID(organizationId, limit, offset)
}

func (service *housingsImpl) FindOneByID(id string) (entities.Housing, error) {
	return service.repository.FindOneByID(id)
}

func (service *housingsImpl) UpdateOneByID(id string, housing entities.Housing) error {
	return service.repository.UpdateOneById(id, housing)
}

func (service *housingsImpl) InactivateOneByID(id string) error {
	data := entities.Housing{
		Status: utils.InactiveStatus,
	}
	return service.repository.UpdateOneById(id, data)
}
