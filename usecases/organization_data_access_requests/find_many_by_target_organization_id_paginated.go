package organization_data_access_requests

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type FindManyByTargetOrganizationIDPaginated interface {
	Execute(actor entities.User, targetOrganizationID string, offset, limit int64) (int64, []entities.OrganizationDataAccessRequest, error)
}

type findManyByTargetOrganizationIDPaginatedImpl struct {
	organizationDataAccessRequestsRepository repositories.OrganizationDataAccessRequests
	organizationsRepository                  repositories.Organizations
}

func NewFindManyByTargetOrganizationIDPaginated(
	organizationDataAccessRequestsRepository repositories.OrganizationDataAccessRequests,
	organizationsRepository repositories.Organizations,
) FindManyByTargetOrganizationIDPaginated {
	return &findManyByTargetOrganizationIDPaginatedImpl{
		organizationDataAccessRequestsRepository: organizationDataAccessRequestsRepository,
		organizationsRepository:                  organizationsRepository,
	}
}

func (uc *findManyByTargetOrganizationIDPaginatedImpl) Execute(actor entities.User, targetOrganizationID string, offset, limit int64) (int64, []entities.OrganizationDataAccessRequest, error) {
	targetOrganization, err := uc.organizationsRepository.FindOneByID(targetOrganizationID)

	if err != nil {
		return 0, nil, err
	}

	if err = guards.IsOrganizationAdmin(actor, targetOrganization); err != nil {
		return 0, nil, err
	}

	return uc.organizationDataAccessRequestsRepository.FindManyByTargetOrganizationIDPaginated(targetOrganizationID, offset, limit)
}
