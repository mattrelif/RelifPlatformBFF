package join_organization_invites

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type FindManyByOrganizationIDPaginated interface {
	Execute(actor entities.User, organizationID string, offset, limit int64) (int64, []entities.JoinOrganizationInvite, error)
}

type findManyByOrganizationIDPaginatedImpl struct {
	joinOrganizationInvitesRepository repositories.JoinOrganizationInvites
	organizationRepository            repositories.Organizations
}

func NewFindManyByOrganizationIDPaginated(
	joinOrganizationInvitesRepository repositories.JoinOrganizationInvites,
	organizationRepository repositories.Organizations,
) FindManyByOrganizationIDPaginated {
	return &findManyByOrganizationIDPaginatedImpl{
		joinOrganizationInvitesRepository: joinOrganizationInvitesRepository,
		organizationRepository:            organizationRepository,
	}
}

func (uc *findManyByOrganizationIDPaginatedImpl) Execute(actor entities.User, organizationID string, offset, limit int64) (int64, []entities.JoinOrganizationInvite, error) {
	organization, err := uc.organizationRepository.FindOneByID(organizationID)

	if err != nil {
		return 0, nil, err
	}

	if err = guards.IsOrganizationAdmin(actor, organization); err != nil {
		return 0, nil, err
	}

	return uc.joinOrganizationInvitesRepository.FindManyByOrganizationIDPaginated(organizationID, offset, limit)
}
