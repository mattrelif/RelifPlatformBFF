package services

import (
	"relif/bff/entities"
	"relif/bff/repositories"
)

type OrganizationDataAccessGrants interface {
	Create(grant entities.OrganizationDataAccessGrant) error
}

type organizationDataAccessGrantsImpl struct {
	repository repositories.OrganizationDataAccessGrants
}

func NewOrganizationDataAccessGrants(repository repositories.OrganizationDataAccessGrants) OrganizationDataAccessGrants {
	return &organizationDataAccessGrantsImpl{
		repository: repository,
	}
}

func (service *organizationDataAccessGrantsImpl) Create(grant entities.OrganizationDataAccessGrant) error {
	return service.repository.Create(grant)
}
