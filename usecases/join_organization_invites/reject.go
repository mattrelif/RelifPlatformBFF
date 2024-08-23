package join_organization_invites

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
	"relif/platform-bff/utils"
	"time"
)

type Reject interface {
	Execute(actor entities.User, inviteID, rejectReason string) error
}

type rejectImpl struct {
	joinOrganizationInvitesRepository repositories.JoinOrganizationInvites
	usersRepository                   repositories.Users
}

func NewReject(
	joinOrganizationInvitesRepository repositories.JoinOrganizationInvites,
	usersRepository repositories.Users,
) Reject {
	return &rejectImpl{
		joinOrganizationInvitesRepository: joinOrganizationInvitesRepository,
		usersRepository:                   usersRepository,
	}
}

func (uc *rejectImpl) Execute(actor entities.User, inviteID, rejectReason string) error {
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

	invite.Status = utils.RejectedStatus
	invite.RejectReason = rejectReason
	invite.RejectedAt = time.Now()

	if err = uc.joinOrganizationInvitesRepository.UpdateOneByID(invite.ID, invite); err != nil {
		return err
	}

	return nil
}
