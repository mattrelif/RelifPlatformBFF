package services

import (
	"relif/bff/entities"
	"relif/bff/repositories"
	"relif/bff/utils"
	"time"
)

type JoinOrganizationRequests interface {
	Create(userId, organizationId string) (entities.JoinOrganizationRequest, error)
	FindManyByOrganizationId(organizationId string, offset, limit int64) (int64, []entities.JoinOrganizationRequest, error)
	FindOneById(id string) (entities.JoinOrganizationRequest, error)
	Accept(id string, auditor entities.User) error
	Reject(id string, auditor entities.User) error
}

type joinOrganizationRequestsImpl struct {
	usersService Users
	repository   repositories.JoinOrganizationRequests
}

func NewJoinOrganizationRequests(usersService Users, repository repositories.JoinOrganizationRequests) JoinOrganizationRequests {
	return &joinOrganizationRequestsImpl{
		usersService: usersService,
		repository:   repository,
	}
}

func (service *joinOrganizationRequestsImpl) Create(userId, organizationId string) (entities.JoinOrganizationRequest, error) {
	data := entities.JoinOrganizationRequest{
		UserID:         userId,
		OrganizationID: organizationId,
	}
	return service.repository.Create(data)
}

func (service *joinOrganizationRequestsImpl) FindManyByOrganizationId(organizationId string, offset, limit int64) (int64, []entities.JoinOrganizationRequest, error) {
	return service.repository.FindManyByOrganizationId(organizationId, offset, limit)
}

func (service *joinOrganizationRequestsImpl) FindOneById(id string) (entities.JoinOrganizationRequest, error) {
	return service.repository.FindOneById(id)
}

func (service *joinOrganizationRequestsImpl) Accept(id string, auditor entities.User) error {
	request, err := service.repository.FindOneById(id)

	if err != nil {
		return err
	}

	if err = service.repository.UpdateOneById(request.ID, entities.JoinOrganizationRequest{AcceptedAt: time.Now(), AuditorID: auditor.ID}); err != nil {
		return err
	}

	data := entities.User{
		OrganizationID: request.OrganizationID,
		PlatformRole:   utils.OrgMemberPlatformRole,
	}

	if err = service.usersService.UpdateOneById(request.UserID, data); err != nil {
		return err
	}

	return nil
}

func (service *joinOrganizationRequestsImpl) Reject(id string, auditor entities.User) error {
	request, err := service.repository.FindOneById(id)

	if err != nil {
		return err
	}

	if err = service.repository.UpdateOneById(request.ID, entities.JoinOrganizationRequest{RejectedAt: time.Now(), AuditorID: auditor.ID}); err != nil {
		return err
	}

	return err
}
