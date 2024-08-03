package services

import (
	"relif/bff/entities"
	"relif/bff/repositories"
	"relif/bff/utils"
)

type Organizations interface {
	Create(data entities.Organization, ownerId string) (entities.Organization, error)
	FindMany(offset, limit int64) (int64, []entities.Organization, error)
	FindOneById(id string) (entities.Organization, error)
	UpdateOneById(id string, data entities.Organization) error
}

type organizationsImpl struct {
	repository   repositories.Organizations
	usersService Users
}

func NewOrganizations(repository repositories.Organizations, usersService Users) Organizations {
	return &organizationsImpl{
		repository:   repository,
		usersService: usersService,
	}
}

func (service *organizationsImpl) Create(data entities.Organization, ownerId string) (entities.Organization, error) {
	data.OwnerID = ownerId
	organization, err := service.repository.Create(data)

	if err != nil {
		return entities.Organization{}, err
	}

	userData := entities.User{
		OrganizationID: organization.ID,
		PlatformRole:   utils.OrgAdminPlatformRole,
	}

	if err = service.usersService.UpdateOneById(ownerId, userData); err != nil {
		return entities.Organization{}, err
	}

	return organization, nil
}

func (service *organizationsImpl) FindMany(offset, limit int64) (int64, []entities.Organization, error) {
	return service.repository.FindMany(offset, limit)
}

func (service *organizationsImpl) FindOneById(organizationId string) (entities.Organization, error) {
	return service.repository.FindOneById(organizationId)
}

func (service *organizationsImpl) UpdateOneById(id string, data entities.Organization) error {
	return service.repository.UpdateOneById(id, data)
}
