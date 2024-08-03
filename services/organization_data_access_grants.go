package services

import (
	"relif/bff/entities"
	"relif/bff/repositories"
	"relif/bff/utils"
)

type OrganizationDataAccessGrants interface {
	Create(grant entities.OrganizationDataAccessGrant) error
	FindManyByOrganizationId(organizationId string, limit, offset int64) (int64, []entities.OrganizationDataAccessGrant, error)
	FindManyByTargetOrganizationId(organizationId string, limit, offset int64) (int64, []entities.OrganizationDataAccessGrant, error)
	FindOneById(id string) (entities.OrganizationDataAccessGrant, error)
	DeleteOneById(id string) error
	ExistsByOrganizationIdAndTargetOrganizationId(organizationId, targetOrganizationId string) (bool, error)
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

func (service *organizationDataAccessGrantsImpl) FindManyByOrganizationId(organizationId string, limit, offset int64) (int64, []entities.OrganizationDataAccessGrant, error) {
	return service.repository.FindManyByOrganizationId(organizationId, limit, offset)
}

func (service *organizationDataAccessGrantsImpl) FindManyByTargetOrganizationId(organizationId string, limit, offset int64) (int64, []entities.OrganizationDataAccessGrant, error) {
	return service.repository.FindManyByTargetOrganizationId(organizationId, limit, offset)
}

func (service *organizationDataAccessGrantsImpl) FindOneById(id string) (entities.OrganizationDataAccessGrant, error) {
	return service.repository.FindOneById(id)
}

func (service *organizationDataAccessGrantsImpl) DeleteOneById(id string) error {
	return service.repository.DeleteOneById(id)
}

func (service *organizationDataAccessGrantsImpl) ExistsByOrganizationIdAndTargetOrganizationId(organizationId, targetOrganizationId string) (bool, error) {
	count, err := service.repository.CountByOrganizationIdAndTargetOrganizationId(organizationId, targetOrganizationId)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (service *organizationDataAccessGrantsImpl) AuthorizeExternalMutation(user entities.User, id string) error {
	grant, err := service.repository.FindOneById(id)

	if err != nil {
		return err
	}

	if user.OrganizationID != grant.TargetOrganizationID && user.PlatformRole != utils.OrgAdminPlatformRole {
		return utils.ErrUnauthorizedAction
	}

	return nil
}
