package services

import (
	"relif/bff/entities"
	"relif/bff/repositories"
	"relif/bff/utils"
)

type Housings interface {
	Create(user entities.User, housing entities.Housing) (entities.Housing, error)
	FindManyByOrganizationID(organizationID, search string, limit, offset int64) (int64, []entities.Housing, error)
	FindOneByID(id string) (entities.Housing, error)
	FindOneCompleteByID(id string) (entities.Housing, error)
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

func (service *housingsImpl) FindManyByOrganizationID(organizationID, search string, limit, offset int64) (int64, []entities.Housing, error) {
	return service.repository.FindManyByOrganizationID(organizationID, search, limit, offset)
}

func (service *housingsImpl) FindOneByID(id string) (entities.Housing, error) {
	return service.repository.FindOneByID(id)
}

func (service *housingsImpl) FindOneCompleteByID(id string) (entities.Housing, error) {
	return service.repository.FindOneCompleteByID(id)
}

func (service *housingsImpl) UpdateOneByID(id string, housing entities.Housing) error {
	return service.repository.UpdateOneByID(id, housing)
}

func (service *housingsImpl) InactivateOneByID(id string) error {
	housing, err := service.FindOneByID(id)

	if err != nil {
		return err
	}

	housing.Status = utils.InactiveStatus

	return service.repository.UpdateOneByID(id, housing)
}
