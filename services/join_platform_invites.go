package services

import (
	"relif/bff/entities"
	"relif/bff/repositories"
	"relif/bff/utils"
)

type JoinPlatformInvites interface {
	Create(data entities.JoinPlatformInvite, inviter entities.User) (entities.JoinPlatformInvite, error)
	AuthorizeAccessManyByOrganizationId(user entities.User, organizationId string) error
	FindManyByOrganizationId(organizationId string, limit, offset int64) (int64, []entities.JoinPlatformInvite, error)
	ConsumeByCode(code string) (entities.JoinPlatformInvite, error)
}

type joinPlatformInvitesImpl struct {
	repository    repositories.JoinPlatformInvites
	emailService  Email
	usersService  Users
	uuidGenerator utils.UuidGenerator
}

func NewJoinPlatformInvites(repository repositories.JoinPlatformInvites, emailService Email, usersService Users, uuidGenerator utils.UuidGenerator) JoinPlatformInvites {
	return &joinPlatformInvitesImpl{
		repository:    repository,
		emailService:  emailService,
		usersService:  usersService,
		uuidGenerator: uuidGenerator,
	}
}

func (service *joinPlatformInvitesImpl) Create(data entities.JoinPlatformInvite, inviter entities.User) (entities.JoinPlatformInvite, error) {
	userExists, err := service.usersService.ExistsByEmail(data.InvitedEmail)

	if err != nil {
		return entities.JoinPlatformInvite{}, err
	}

	if userExists {
		return entities.JoinPlatformInvite{}, utils.ErrUserAlreadyExists
	}

	data.Code = service.uuidGenerator()
	data.InviterID = inviter.ID
	data.OrganizationID = inviter.OrganizationID

	invite, err := service.repository.Create(data)

	if err != nil {
		return entities.JoinPlatformInvite{}, err
	}

	if err = service.emailService.SendPlatformInviteEmail(inviter, invite.InvitedEmail, invite.Code); err != nil {
		return entities.JoinPlatformInvite{}, err
	}

	return invite, nil
}

func (service *joinPlatformInvitesImpl) FindManyByOrganizationId(organizationId string, limit, offset int64) (int64, []entities.JoinPlatformInvite, error) {
	return service.repository.FindManyByOrganizationId(organizationId, limit, offset)
}

func (service *joinPlatformInvitesImpl) ConsumeByCode(code string) (entities.JoinPlatformInvite, error) {
	invite, err := service.repository.FindOneAndDeleteByCode(code)

	if err != nil {
		return entities.JoinPlatformInvite{}, err
	}

	return invite, nil
}

func (service *joinPlatformInvitesImpl) AuthorizeAccessManyByOrganizationId(user entities.User, organizationId string) error {
	if (user.OrganizationID != organizationId && user.PlatformRole != utils.OrgAdminPlatformRole) && user.PlatformRole != utils.RelifMemberPlatformRole {
		return utils.ErrUnauthorizedAction
	}

	return nil
}
