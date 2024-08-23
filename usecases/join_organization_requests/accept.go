package join_organization_requests

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
	"relif/platform-bff/utils"
	"time"
)

type Accept interface {
	Execute(actor entities.User, requestID string) error
}

type acceptImpl struct {
	joinOrganizationRequestsRepository repositories.JoinOrganizationRequests
	usersRepository                    repositories.Users
	organizationsRepository            repositories.Organizations
}

func NewAccept(
	joinOrganizationRequestsRepository repositories.JoinOrganizationRequests,
	usersRepository repositories.Users,
	organizationsRepository repositories.Organizations,
) Accept {
	return &acceptImpl{
		joinOrganizationRequestsRepository: joinOrganizationRequestsRepository,
		usersRepository:                    usersRepository,
		organizationsRepository:            organizationsRepository,
	}
}

func (uc *acceptImpl) Execute(actor entities.User, requestID string) error {
	request, err := uc.joinOrganizationRequestsRepository.FindOneByID(requestID)

	if err != nil {
		return err
	}

	organization, err := uc.organizationsRepository.FindOneByID(request.OrganizationID)

	if err != nil {
		return err
	}

	if err = guards.IsOrganizationAdmin(actor, organization); err != nil {
		return err
	}

	requester, err := uc.usersRepository.FindOneByID(request.UserID)

	if err != nil {
		return err
	}

	requester.OrganizationID = request.OrganizationID
	requester.PlatformRole = utils.OrgMemberPlatformRole

	if err = uc.usersRepository.UpdateOneByID(requester.ID, requester); err != nil {
		return err
	}

	request.Status = utils.AcceptedStatus
	request.AcceptedAt = time.Now()
	request.AuditorID = actor.ID

	if err = uc.joinOrganizationRequestsRepository.UpdateOneByID(request.ID, request); err != nil {
		return err
	}

	return nil
}
