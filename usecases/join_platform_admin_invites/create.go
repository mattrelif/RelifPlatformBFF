package join_platform_admin_invites

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
	"relif/platform-bff/services"
	"relif/platform-bff/utils"
)

type Create interface {
	Execute(actor entities.User, data entities.JoinPlatformAdminInvite) (entities.JoinPlatformAdminInvite, error)
}

type createImpl struct {
	repository            repositories.JoinPlatformAdminInvites
	emailService          services.Email
	uuidGeneratorFunction utils.UuidGenerator
}

func NewCreate(
	repository repositories.JoinPlatformAdminInvites,
	emailService services.Email,
	uuidGeneratorFunction utils.UuidGenerator,
) Create {
	return &createImpl{
		repository:            repository,
		emailService:          emailService,
		uuidGeneratorFunction: uuidGeneratorFunction,
	}
}

func (uc *createImpl) Execute(actor entities.User, data entities.JoinPlatformAdminInvite) (entities.JoinPlatformAdminInvite, error) {
	if err := guards.IsSuperUser(actor); err != nil {
		return entities.JoinPlatformAdminInvite{}, err
	}

	data.Code = uc.uuidGeneratorFunction()
	data.InviterID = actor.ID

	invite, err := uc.repository.Create(data)

	if err != nil {
		return entities.JoinPlatformAdminInvite{}, err
	}

	if err = uc.emailService.SendPlatformAdminInviteEmail(actor, invite); err != nil {
		return entities.JoinPlatformAdminInvite{}, err
	}

	return invite, nil
}
