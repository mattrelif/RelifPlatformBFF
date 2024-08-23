package join_organization_requests

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type FindManyByOrganizationIDPaginated interface {
	Execute(actor entities.User, organizationID string, offset, limit int64) (int64, []entities.JoinOrganizationRequest, error)
}

type findManyByOrganizationIDPaginatedImpl struct {
	joinOrganizationInvitesRepository repositories.JoinOrganizationRequests
	organizationRepository            repositories.Organizations
}

func NewFindManyByOrganizationIDPaginated(
	joinOrganizationInvitesRepository repositories.JoinOrganizationRequests,
	organizationRepository repositories.Organizations,
) FindManyByOrganizationIDPaginated {
	return &findManyByOrganizationIDPaginatedImpl{
		joinOrganizationInvitesRepository: joinOrganizationInvitesRepository,
		organizationRepository:            organizationRepository,
	}
}

func (uc *findManyByOrganizationIDPaginatedImpl) Execute(actor entities.User, organizationID string, offset, limit int64) (int64, []entities.JoinOrganizationRequest, error) {
	organization, err := uc.organizationRepository.FindOneByID(organizationID)

	if err != nil {
		return 0, nil, err
	}

	if err = guards.IsOrganizationAdmin(actor, organization); err != nil {
		return 0, nil, err
	}

	return uc.joinOrganizationInvitesRepository.FindManyByOrganizationIDPaginated(organizationID, offset, limit)
}
