package join_organization_invites

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type Create interface {
	Execute(actor entities.User, organizationID string, data entities.JoinOrganizationInvite) (entities.JoinOrganizationInvite, error)
}

type createImpl struct {
	joinOrganizationInvitesRepository repositories.JoinOrganizationInvites
	organizationsRepository           repositories.Organizations
}

func NewCreate(
	joinOrganizationInvitesRepository repositories.JoinOrganizationInvites,
	organizationsRepository repositories.Organizations,
) Create {
	return &createImpl{
		joinOrganizationInvitesRepository: joinOrganizationInvitesRepository,
		organizationsRepository:           organizationsRepository,
	}
}

func (uc *createImpl) Execute(actor entities.User, organizationID string, data entities.JoinOrganizationInvite) (entities.JoinOrganizationInvite, error) {
	organization, err := uc.organizationsRepository.FindOneByID(organizationID)

	if err != nil {
		return entities.JoinOrganizationInvite{}, err
	}

	if err = guards.IsOrganizationAdmin(actor, organization); err != nil {
		return entities.JoinOrganizationInvite{}, err
	}

	data.OrganizationID = organization.ID
	data.CreatorID = actor.ID

	return uc.joinOrganizationInvitesRepository.Create(data)
}
