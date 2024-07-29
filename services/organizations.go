package services

import (
	"relif/bff/entities"
	"relif/bff/repositories"
)

type Organizations interface {
	Create(data entities.Organization, creatorId string) (entities.Organization, error)
	FindMany(offset, limit int64) (int64, []entities.Organization, error)
	FindOneAndUpdateById(id string, data entities.Organization) (entities.Organization, error)
	UpdateOneById(id string, data entities.Organization) error
}

type organizationsImpl struct {
	repository repositories.Organizations
}

func NewOrganizations(repository repositories.Organizations) Organizations {
	return &organizationsImpl{
		repository: repository,
	}
}

func (service *organizationsImpl) Create(organization entities.Organization, creatorId string) (entities.Organization, error) {
	organization.CreatorID = creatorId
	return service.repository.Create(organization)
}

func (service *organizationsImpl) FindMany(offset, limit int64) (int64, []entities.Organization, error) {
	return service.repository.FindMany(offset, limit)
}

func (service *organizationsImpl) FindOneAndUpdateById(id string, data entities.Organization) (entities.Organization, error) {
	return service.repository.FindOneAndUpdateById(id, data)
}

func (service *organizationsImpl) UpdateOneById(id string, data entities.Organization) error {
	return service.repository.UpdateOneById(id, data)
}
