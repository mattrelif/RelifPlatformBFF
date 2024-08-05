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
	FindManyByUserId(userId string, offset, limit int64) (int64, []entities.JoinOrganizationRequest, error)
	FindOneById(id string) (entities.JoinOrganizationRequest, error)
	Accept(id string, auditor entities.User) error
	Reject(id string, auditor entities.User, data entities.JoinOrganizationRequest) error
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

func (service *joinOrganizationRequestsImpl) FindManyByUserId(userId string, offset, limit int64) (int64, []entities.JoinOrganizationRequest, error) {
	return service.repository.FindManyByUserId(userId, offset, limit)
}

func (service *joinOrganizationRequestsImpl) FindOneById(id string) (entities.JoinOrganizationRequest, error) {
	return service.repository.FindOneById(id)
}

func (service *joinOrganizationRequestsImpl) Accept(id string, auditor entities.User) error {
	request, err := service.repository.FindOneById(id)

	if err != nil {
		return err
	}

	user, err := service.usersService.FindOneById(request.UserID)

	if err != nil {
		return err
	}

	request.AcceptedAt = time.Now()
	request.Status = utils.AcceptedStatus
	request.AuditorID = auditor.ID

	if err = service.repository.UpdateOneById(request.ID, request); err != nil {
		return err
	}

	user.OrganizationID = request.OrganizationID
	user.PlatformRole = utils.OrgMemberPlatformRole

	if err = service.usersService.UpdateOneById(user.ID, user); err != nil {
		return err
	}

	return nil
}

func (service *joinOrganizationRequestsImpl) Reject(id string, auditor entities.User, data entities.JoinOrganizationRequest) error {
	request, err := service.repository.FindOneById(id)

	if err != nil {
		return err
	}

	request.AuditorID = auditor.ID
	request.Status = utils.RejectedStatus
	request.RejectedAt = time.Now()
	request.RejectReason = data.RejectReason

	if err = service.repository.UpdateOneById(request.ID, request); err != nil {
		return err
	}

	return err
}
