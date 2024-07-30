package services

import (
	"relif/bff/entities"
	"relif/bff/repositories"
)

type VoluntaryPeople interface {
	Create(user entities.User, data entities.VoluntaryPerson) (entities.VoluntaryPerson, error)
	FindManyByOrganizationId(organizationId string, limit, offset int64) (int64, []entities.VoluntaryPerson, error)
	FindOneById(id string) (entities.VoluntaryPerson, error)
	FindOneAndUpdateById(id string, data entities.VoluntaryPerson) (entities.VoluntaryPerson, error)
	DeleteOneById(id string) error
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
	data.Status = "ACTIVE"
	data.OrganizationID = user.OrganizationID
	return service.repository.Create(data)
}

func (service *voluntaryPeopleImpl) FindManyByOrganizationId(organizationId string, limit, offset int64) (int64, []entities.VoluntaryPerson, error) {
	return service.repository.FindManyByOrganizationId(organizationId, limit, offset)
}

func (service *voluntaryPeopleImpl) FindOneById(id string) (entities.VoluntaryPerson, error) {
	return service.repository.FindOneById(id)
}

func (service *voluntaryPeopleImpl) FindOneAndUpdateById(id string, data entities.VoluntaryPerson) (entities.VoluntaryPerson, error) {
	return service.repository.FindOneAndUpdateById(id, data)
}

func (service *voluntaryPeopleImpl) DeleteOneById(id string) error {
	return service.repository.DeleteOneById(id)
}
