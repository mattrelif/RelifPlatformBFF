package services

import (
	"relif/bff/entities"
	"relif/bff/repositories"
	"time"
)

type Organizations interface {
	Create(organization entities.Organization, creatorId string) (string, error)
	FindMany(offset, limit int64) (int64, []entities.Organization, error)
	FindOneAndUpdateById(id string, organization entities.Organization) (entities.Organization, error)
	UpdateOneById(id string, organization entities.Organization) error
}

type organizationsImpl struct {
	repository repositories.Organizations
}

func NewOrganizations(repository repositories.Organizations) Organizations {
	return &organizationsImpl{
		repository: repository,
	}
}

func (service *organizationsImpl) Create(organization entities.Organization, creatorId string) (string, error) {
	organization.CreatorID = creatorId
	organization.CreatedAt = time.Now()

	return service.repository.Create(organization)
}

func (service *organizationsImpl) FindMany(offset, limit int64) (int64, []entities.Organization, error) {
	return service.repository.FindMany(offset, limit)
}

func (service *organizationsImpl) FindOneAndUpdateById(id string, organization entities.Organization) (entities.Organization, error) {
	organization.UpdatedAt = time.Now()
	return service.repository.FindOneAndUpdateById(id, organization)
}

func (service *organizationsImpl) UpdateOneById(id string, organization entities.Organization) error {
	organization.UpdatedAt = time.Now()
	return service.repository.UpdateOneById(id, organization)
}
