package join_platform_invites

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type FindManyByOrganizationIDPaginated interface {
	Execute(actor entities.User, organizationID string, offset, limit int64) (int64, []entities.JoinPlatformInvite, error)
}

type findManyByOrganizationIDPaginatedImpl struct {
	joinPlatformInvitesRepository repositories.JoinPlatformInvites
	organizationsRepository       repositories.Organizations
}

func NewFindManyByOrganizationIDPaginated(
	joinPlatformInvitesRepository repositories.JoinPlatformInvites,
	organizationsRepository repositories.Organizations,
) FindManyByOrganizationIDPaginated {
	return &findManyByOrganizationIDPaginatedImpl{
		joinPlatformInvitesRepository: joinPlatformInvitesRepository,
		organizationsRepository:       organizationsRepository,
	}
}

func (uc *findManyByOrganizationIDPaginatedImpl) Execute(actor entities.User, organizationID string, offset, limit int64) (int64, []entities.JoinPlatformInvite, error) {
	organization, err := uc.organizationsRepository.FindOneByID(organizationID)

	if err != nil {
		return 0, nil, err
	}

	if err = guards.IsOrganizationAdmin(actor, organization); err != nil {
		return 0, nil, err
	}

	return uc.joinPlatformInvitesRepository.FindManyByOrganizationIDPaginated(organization.ID, offset, limit)
}
