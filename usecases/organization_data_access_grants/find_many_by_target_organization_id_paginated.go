package organization_data_access_grants

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type FindManyByTargetOrganizationIDPaginated interface {
	Execute(actor entities.User, targetOrganizationID string, offset, limit int64) (int64, []entities.OrganizationDataAccessGrant, error)
}

type findManyByTargetOrganizationIDPaginatedImpl struct {
	organizationDataAccessGrantsRepository repositories.OrganizationDataAccessGrants
	organizationsRepository                repositories.Organizations
}

func NewFindManyByTargetOrganizationIDPaginated(
	organizationDataAccessGrantsRepository repositories.OrganizationDataAccessGrants,
	organizationsRepository repositories.Organizations,
) FindManyByTargetOrganizationIDPaginated {
	return &findManyByTargetOrganizationIDPaginatedImpl{
		organizationDataAccessGrantsRepository: organizationDataAccessGrantsRepository,
		organizationsRepository:                organizationsRepository,
	}
}

func (uc *findManyByTargetOrganizationIDPaginatedImpl) Execute(actor entities.User, targetOrganizationID string, offset, limit int64) (int64, []entities.OrganizationDataAccessGrant, error) {
	targetOrganization, err := uc.organizationsRepository.FindOneByID(targetOrganizationID)

	if err != nil {
		return 0, nil, err
	}

	if err = guards.IsOrganizationAdmin(actor, targetOrganization); err != nil {
		return 0, nil, err
	}

	return uc.organizationDataAccessGrantsRepository.FindManyByTargetOrganizationIDPaginated(targetOrganization.ID, offset, limit)
}
