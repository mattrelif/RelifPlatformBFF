package services

import (
	"relif/bff/entities"
	"relif/bff/repositories"
)

type OrganizationDataAccessGrants interface {
	Create(grant entities.OrganizationDataAccessGrant) error
	FindManyByOrganizationID(organizationID string, limit, offset int64) (int64, []entities.OrganizationDataAccessGrant, error)
	FindManyByTargetOrganizationID(organizationID string, limit, offset int64) (int64, []entities.OrganizationDataAccessGrant, error)
	FindOneByID(id string) (entities.OrganizationDataAccessGrant, error)
	DeleteOneByID(id string) error
	ExistsByOrganizationIDAndTargetOrganizationID(organizationID, targetOrganizationID string) (bool, error)
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

func (service *organizationDataAccessGrantsImpl) FindManyByOrganizationID(organizationID string, limit, offset int64) (int64, []entities.OrganizationDataAccessGrant, error) {
	return service.repository.FindManyByOrganizationID(organizationID, limit, offset)
}

func (service *organizationDataAccessGrantsImpl) FindManyByTargetOrganizationID(organizationID string, limit, offset int64) (int64, []entities.OrganizationDataAccessGrant, error) {
	return service.repository.FindManyByTargetOrganizationID(organizationID, limit, offset)
}

func (service *organizationDataAccessGrantsImpl) FindOneByID(id string) (entities.OrganizationDataAccessGrant, error) {
	return service.repository.FindOneByID(id)
}

func (service *organizationDataAccessGrantsImpl) DeleteOneByID(id string) error {
	return service.repository.DeleteOneByID(id)
}

func (service *organizationDataAccessGrantsImpl) ExistsByOrganizationIDAndTargetOrganizationID(organizationID, targetOrganizationID string) (bool, error) {
	count, err := service.repository.CountByOrganizationIDAndTargetOrganizationID(organizationID, targetOrganizationID)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
