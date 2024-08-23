package organization_data_access_requests

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type FindManyByRequesterOrganizationIDPaginated interface {
	Execute(actor entities.User, requesterOrganizationID string, offset, limit int64) (int64, []entities.OrganizationDataAccessRequest, error)
}

type findManyByRequesterOrganizationIDPaginatedImpl struct {
	organizationDataAccessRequestsRepository repositories.OrganizationDataAccessRequests
	organizationsRepository                  repositories.Organizations
}

func NewFindManyByRequesterOrganizationIDPaginated(
	organizationDataAccessRequestsRepository repositories.OrganizationDataAccessRequests,
	organizationsRepository repositories.Organizations,
) FindManyByRequesterOrganizationIDPaginated {
	return &findManyByRequesterOrganizationIDPaginatedImpl{
		organizationDataAccessRequestsRepository: organizationDataAccessRequestsRepository,
		organizationsRepository:                  organizationsRepository,
	}
}

func (uc *findManyByRequesterOrganizationIDPaginatedImpl) Execute(actor entities.User, requesterOrganizationID string, offset, limit int64) (int64, []entities.OrganizationDataAccessRequest, error) {
	requesterOrganization, err := uc.organizationsRepository.FindOneByID(requesterOrganizationID)

	if err != nil {
		return 0, nil, err
	}

	if err = guards.IsOrganizationAdmin(actor, requesterOrganization); err != nil {
		return 0, nil, err
	}

	return uc.organizationDataAccessRequestsRepository.FindManyByRequesterOrganizationIDPaginated(requesterOrganizationID, offset, limit)
}
