package join_organization_invites

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
	"relif/platform-bff/utils"
	"time"
)

type Accept interface {
	Execute(actor entities.User, inviteID string) error
}

type acceptImpl struct {
	joinOrganizationInvitesRepository repositories.JoinOrganizationInvites
	usersRepository                   repositories.Users
}

func NewAccept(
	joinOrganizationInvitesRepository repositories.JoinOrganizationInvites,
	usersRepository repositories.Users,
) Accept {
	return &acceptImpl{
		joinOrganizationInvitesRepository: joinOrganizationInvitesRepository,
		usersRepository:                   usersRepository,
	}
}

func (uc *acceptImpl) Execute(actor entities.User, inviteID string) error {
	invite, err := uc.joinOrganizationInvitesRepository.FindOneByID(inviteID)

	if err != nil {
		return err
	}

	invited, err := uc.usersRepository.FindOneByID(invite.UserID)

	if err != nil {
		return err
	}

	if err = guards.IsUser(actor, invited); err != nil {
		return err
	}

	invited.OrganizationID = invite.OrganizationID
	invited.PlatformRole = utils.OrgMemberPlatformRole

	if err = uc.usersRepository.UpdateOneByID(invited.ID, invited); err != nil {
		return err
	}

	invite.Status = utils.AcceptedStatus
	invite.AcceptedAt = time.Now()

	if err = uc.joinOrganizationInvitesRepository.UpdateOneByID(invite.ID, invite); err != nil {
		return err
	}

	return nil
}
