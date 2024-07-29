package services

import (
	"relif/bff/entities"
	"relif/bff/repositories"
)

type JoinOrganizationInvites interface {
	Create(userID string, data entities.JoinOrganizationInvite) (entities.JoinOrganizationInvite, error)
	FindManyByOrganizationId(organizationId string, offset, limit int64) (int64, []entities.JoinOrganizationInvite, error)
	Accept(id string) (string, error)
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

func (service *joinOrganizationInvitesImpl) Create(userID string, data entities.JoinOrganizationInvite) (entities.JoinOrganizationInvite, error) {
	data.CreatorID = userID
	return service.repository.Create(data)
}

func (service *joinOrganizationInvitesImpl) FindManyByOrganizationId(organizationId string, offset, limit int64) (int64, []entities.JoinOrganizationInvite, error) {
	return service.repository.FindManyByOrganizationId(organizationId, offset, limit)
}

func (service *joinOrganizationInvitesImpl) Accept(id string) (string, error) {
	invite, err := service.repository.FindOneAndDeleteById(id)

	if err != nil {
		return "", err
	}

	data := entities.User{
		OrganizationID: invite.OrganizationID,
	}

	if err = service.usersService.UpdateOneById(invite.UserID, data); err != nil {
		return "", err
	}

	return invite.OrganizationID, nil
}

func (service *joinOrganizationInvitesImpl) Reject(id string) error {
	return service.repository.DeleteOneById(id)
}
