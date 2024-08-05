package services

import (
	"relif/bff/entities"
	"relif/bff/repositories"
	"relif/bff/utils"
)

type VoluntaryPeople interface {
	Create(user entities.User, data entities.VoluntaryPerson) (entities.VoluntaryPerson, error)
	FindManyByOrganizationId(organizationId, search string, limit, offset int64) (int64, []entities.VoluntaryPerson, error)
	FindOneById(id string) (entities.VoluntaryPerson, error)
	UpdateOneById(id string, data entities.VoluntaryPerson) error
	InactivateOneById(id string) error
}

type voluntaryPeopleImpl struct {
	repository repositories.VoluntaryPeople
}

func NewVoluntaryPeople(repository repositories.VoluntaryPeople) VoluntaryPeople {
	return &voluntaryPeopleImpl{
		repository: repository,
	}
}

func (service *voluntaryPeopleImpl) Create(user entities.User, data entities.VoluntaryPerson) (entities.VoluntaryPerson, error) {
	exists, err := service.ExistsOneByEmail(data.Email)

	if err != nil {
		return entities.VoluntaryPerson{}, err
	}

	if exists {
		return entities.VoluntaryPerson{}, utils.ErrVoluntaryPersonAlreadyExists
	}

	data.OrganizationID = user.OrganizationID

	return service.repository.Create(data)
}

func (service *voluntaryPeopleImpl) FindManyByOrganizationId(organizationId, search string, limit, offset int64) (int64, []entities.VoluntaryPerson, error) {
	return service.repository.FindManyByOrganizationId(organizationId, search, limit, offset)
}

func (service *voluntaryPeopleImpl) FindOneById(id string) (entities.VoluntaryPerson, error) {
	return service.repository.FindOneById(id)
}

func (service *voluntaryPeopleImpl) UpdateOneById(id string, data entities.VoluntaryPerson) error {
	return service.repository.UpdateOneById(id, data)
}

func (service *voluntaryPeopleImpl) InactivateOneById(id string) error {
	voluntary, err := service.FindOneById(id)

	if err != nil {
		return err
	}

	voluntary.OrganizationID = utils.InactiveStatus

	return service.repository.UpdateOneById(id, voluntary)
}

func (service *voluntaryPeopleImpl) ExistsOneByEmail(email string) (bool, error) {
	count, err := service.repository.CountByEmail(email)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
