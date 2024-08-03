package services

import (
	"relif/bff/entities"
	"relif/bff/repositories"
	"relif/bff/utils"
	"time"
)

type JoinOrganizationInvites interface {
	Create(user entities.User, data entities.JoinOrganizationInvite) (entities.JoinOrganizationInvite, error)
	FindManyByOrganizationId(organizationId string, offset, limit int64) (int64, []entities.JoinOrganizationInvite, error)
	FindManyByUserId(userId string, offset, limit int64) (int64, []entities.JoinOrganizationInvite, error)
	FindOneById(id string) (entities.JoinOrganizationInvite, error)
	Accept(id string) error
	Reject(id string) error
}

type joinOrganizationInvitesImpl struct {
	usersService Users
	repository   repositories.JoinOrganizationInvites
}

func NewJoinOrganizationInvites(usersService Users, repository repositories.JoinOrganizationInvites) JoinOrganizationInvites {
	return &joinOrganizationInvitesImpl{
		usersService: usersService,
		repository:   repository,
	}
}

func (service *joinOrganizationInvitesImpl) Create(user entities.User, data entities.JoinOrganizationInvite) (entities.JoinOrganizationInvite, error) {
	data.CreatorID = user.ID
	data.OrganizationID = user.OrganizationID
	return service.repository.Create(data)
}

func (service *joinOrganizationInvitesImpl) FindManyByOrganizationId(organizationId string, offset, limit int64) (int64, []entities.JoinOrganizationInvite, error) {
	return service.repository.FindManyByOrganizationId(organizationId, offset, limit)
}

func (service *joinOrganizationInvitesImpl) FindManyByUserId(userId string, offset, limit int64) (int64, []entities.JoinOrganizationInvite, error) {
	return service.repository.FindManyByUserId(userId, offset, limit)
}

func (service *joinOrganizationInvitesImpl) FindOneById(id string) (entities.JoinOrganizationInvite, error) {
	return service.repository.FindOneById(id)
}

func (service *joinOrganizationInvitesImpl) Accept(id string) error {
	invite, err := service.repository.FindOneById(id)

	if err != nil {
		return err
	}

	if err = service.repository.UpdateOneById(invite.ID, entities.JoinOrganizationInvite{AcceptedAt: time.Now(), Status: utils.AcceptedStatus}); err != nil {
		return err
	}

	data := entities.User{
		OrganizationID: invite.OrganizationID,
		PlatformRole:   utils.OrgMemberPlatformRole,
	}

	if err = service.usersService.UpdateOneById(invite.UserID, data); err != nil {
		return err
	}

	return nil
}

func (service *joinOrganizationInvitesImpl) Reject(id string) error {
	invite, err := service.repository.FindOneById(id)

	if err != nil {
		return err
	}

	if err = service.repository.UpdateOneById(invite.ID, entities.JoinOrganizationInvite{RejectedAt: time.Now(), Status: utils.RejectedStatus}); err != nil {
		return err
	}

	return nil
}
