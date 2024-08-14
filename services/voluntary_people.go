package services

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/repositories"
	"relif/platform-bff/utils"
)

type VoluntaryPeople interface {
	Create(organizationID string, data entities.VoluntaryPerson) (entities.VoluntaryPerson, error)
	FindManyByOrganizationID(organizationID, search string, limit, offset int64) (int64, []entities.VoluntaryPerson, error)
	FindOneByID(id string) (entities.VoluntaryPerson, error)
	UpdateOneByID(id string, data entities.VoluntaryPerson) error
	InactivateOneByID(id string) error
}

type voluntaryPeopleImpl struct {
	repository repositories.VoluntaryPeople
}

func NewVoluntaryPeople(repository repositories.VoluntaryPeople) VoluntaryPeople {
	return &voluntaryPeopleImpl{
		repository: repository,
	}
}

func (service *voluntaryPeopleImpl) Create(organizationID string, data entities.VoluntaryPerson) (entities.VoluntaryPerson, error) {
	exists, err := service.ExistsOneByEmail(data.Email)

	if err != nil {
		return entities.VoluntaryPerson{}, err
	}

	if exists {
		return entities.VoluntaryPerson{}, utils.ErrVoluntaryPersonAlreadyExists
	}

	data.OrganizationID = organizationID

	return service.repository.Create(data)
}

func (service *voluntaryPeopleImpl) FindManyByOrganizationID(organizationID, search string, limit, offset int64) (int64, []entities.VoluntaryPerson, error) {
	return service.repository.FindManyByOrganizationID(organizationID, search, limit, offset)
}

func (service *voluntaryPeopleImpl) FindOneByID(id string) (entities.VoluntaryPerson, error) {
	return service.repository.FindOneByID(id)
}

func (service *voluntaryPeopleImpl) UpdateOneByID(id string, data entities.VoluntaryPerson) error {
	return service.repository.UpdateOneByID(id, data)
}

func (service *voluntaryPeopleImpl) InactivateOneByID(id string) error {
	voluntary, err := service.FindOneByID(id)

	if err != nil {
		return err
	}

	voluntary.Status = utils.InactiveStatus

	return service.repository.UpdateOneByID(id, voluntary)
}

func (service *voluntaryPeopleImpl) ExistsOneByEmail(email string) (bool, error) {
	count, err := service.repository.CountByEmail(email)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
