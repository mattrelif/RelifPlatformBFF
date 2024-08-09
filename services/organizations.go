package services

import (
	"relif/bff/entities"
	"relif/bff/repositories"
	"relif/bff/utils"
)

type Organizations interface {
	Create(data entities.Organization, owner entities.User) (entities.Organization, error)
	FindMany(offset, limit int64) (int64, []entities.Organization, error)
	FindOneByID(id string) (entities.Organization, error)
	UpdateOneByID(id string, data entities.Organization) error
	InactivateOneByID(id string) error
	ReactivateOneByID(id string) error
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

func (service *organizationsImpl) Create(data entities.Organization, owner entities.User) (entities.Organization, error) {
	data.OwnerID = owner.ID

	organization, err := service.repository.Create(data)

	if err != nil {
		return entities.Organization{}, err
	}

	owner.OrganizationID = organization.ID
	owner.PlatformRole = utils.OrgAdminPlatformRole

	if err = service.usersService.UpdateOneByID(owner.ID, owner); err != nil {
		return entities.Organization{}, err
	}

	return organization, nil
}

func (service *organizationsImpl) FindMany(offset, limit int64) (int64, []entities.Organization, error) {
	return service.repository.FindMany(offset, limit)
}

func (service *organizationsImpl) FindOneByID(organizationID string) (entities.Organization, error) {
	return service.repository.FindOneByID(organizationID)
}

func (service *organizationsImpl) UpdateOneByID(id string, data entities.Organization) error {
	return service.repository.UpdateOneByID(id, data)
}

func (service *organizationsImpl) InactivateOneByID(id string) error {
	organization, err := service.FindOneByID(id)

	if err != nil {
		return err
	}

	organization.Status = utils.InactiveStatus

	return service.repository.UpdateOneByID(organization.ID, organization)
}

func (service *organizationsImpl) ReactivateOneByID(id string) error {
	organization, err := service.FindOneByID(id)

	if err != nil {
		return err
	}

	organization.Status = utils.ActiveStatus

	return service.repository.UpdateOneByID(organization.ID, organization)
}
