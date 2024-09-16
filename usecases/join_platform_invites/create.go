package join_platform_invites

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
	"relif/platform-bff/services"
	"relif/platform-bff/utils"
)

type Create interface {
	Execute(actor entities.User, data entities.JoinPlatformInvite) (entities.JoinPlatformInvite, error)
}

type createImpl struct {
	joinPlatformInvitesRepository repositories.JoinPlatformInvites
	organizationsRepository       repositories.Organizations
	usersRepository               repositories.Users
	emailService                  services.Email
	uuidGeneratorFunction         utils.UuidGenerator
}

func NewCreate(
	joinPlatformInvitesRepository repositories.JoinPlatformInvites,
	organizationsRepository repositories.Organizations,
	usersRepository repositories.Users,
	emailService services.Email,
	uuidGeneratorFunction utils.UuidGenerator,
) Create {
	return &createImpl{
		joinPlatformInvitesRepository: joinPlatformInvitesRepository,
		organizationsRepository:       organizationsRepository,
		usersRepository:               usersRepository,
		emailService:                  emailService,
		uuidGeneratorFunction:         uuidGeneratorFunction,
	}
}

func (uc *createImpl) Execute(actor entities.User, data entities.JoinPlatformInvite) (entities.JoinPlatformInvite, error) {
	if err := guards.IsOrganizationAdmin(actor, actor.Organization); err != nil {
		return entities.JoinPlatformInvite{}, err
	}

	count, err := uc.joinPlatformInvitesRepository.CountByInvitedEmail(data.InvitedEmail)

	if err != nil {
		return entities.JoinPlatformInvite{}, err
	}

	if count > 0 {
		return entities.JoinPlatformInvite{}, utils.ErrInviteAlreadyExists
	}

	data.OrganizationID = actor.Organization.ID
	data.InviterID = actor.ID
	data.Code = uc.uuidGeneratorFunction()

	invite, err := uc.joinPlatformInvitesRepository.Create(data)

	if err != nil {
		return entities.JoinPlatformInvite{}, err
	}

	if err = uc.emailService.SendPlatformInviteEmail(actor, invite); err != nil {
		return entities.JoinPlatformInvite{}, err
	}

	return invite, nil
}
