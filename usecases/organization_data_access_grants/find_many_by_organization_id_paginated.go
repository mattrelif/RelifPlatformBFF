package organization_data_access_grants

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type FindManyByOrganizationIDPaginated interface {
	Execute(actor entities.User, organizationID string, offset, limit int64) (int64, []entities.OrganizationDataAccessGrant, error)
}

type findManyByOrganizationIDPaginatedImpl struct {
	organizationDataAccessGrantsRepository repositories.OrganizationDataAccessGrants
	organizationsRepository                repositories.Organizations
}

func NewFindManyByOrganizationIDPaginated(
	organizationDataAccessGrantsRepository repositories.OrganizationDataAccessGrants,
	organizationsRepository repositories.Organizations,
) FindManyByOrganizationIDPaginated {
	return &findManyByOrganizationIDPaginatedImpl{
		organizationDataAccessGrantsRepository: organizationDataAccessGrantsRepository,
		organizationsRepository:                organizationsRepository,
	}
}

func (uc *findManyByOrganizationIDPaginatedImpl) Execute(actor entities.User, organizationID string, offset, limit int64) (int64, []entities.OrganizationDataAccessGrant, error) {
	organization, err := uc.organizationsRepository.FindOneByID(organizationID)

	if err != nil {
		return 0, nil, err
	}

	if err = guards.IsOrganizationAdmin(actor, organization); err != nil {
		return 0, nil, err
	}

	return uc.organizationDataAccessGrantsRepository.FindManyByOrganizationIDPaginated(organization.ID, offset, limit)
}
