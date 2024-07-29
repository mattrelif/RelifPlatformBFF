package services

import (
	"relif/bff/entities"
	"relif/bff/repositories"
	"time"
)

type Housings interface {
	Create(user entities.User, housing entities.Housing) (entities.Housing, error)
	FindManyByOrganizationID(organizationId string, limit, offset int64) (int64, []entities.Housing, error)
	FindOneByID(id string) (entities.Housing, error)
	FindOneAndUpdateById(id string, housing entities.Housing) (entities.Housing, error)
	DeleteOneById(id string) error
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
	housing.CreatedAt = time.Now()
	housing.OrganizationID = user.OrganizationID

	return service.repository.Create(housing)
}

func (service *housingsImpl) FindManyByOrganizationID(organizationId string, limit, offset int64) (int64, []entities.Housing, error) {
	return service.repository.FindManyByOrganizationID(organizationId, limit, offset)
}

func (service *housingsImpl) FindOneByID(id string) (entities.Housing, error) {
	return service.repository.FindOneByID(id)
}

func (service *housingsImpl) FindOneAndUpdateById(id string, housing entities.Housing) (entities.Housing, error) {
	housing.UpdatedAt = time.Now()

	return service.repository.FindOneAndUpdateById(id, housing)
}

func (service *housingsImpl) DeleteOneById(id string) error {
	return service.repository.DeleteOneById(id)
}
