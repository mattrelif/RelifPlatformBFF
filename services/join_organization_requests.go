package services

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/repositories"
	"relif/platform-bff/utils"
	"time"
)

type JoinOrganizationRequests interface {
	Create(userID, organizationID string) (entities.JoinOrganizationRequest, error)
	FindManyByOrganizationID(organizationID string, offset, limit int64) (int64, []entities.JoinOrganizationRequest, error)
	FindManyByUserID(userID string, offset, limit int64) (int64, []entities.JoinOrganizationRequest, error)
	FindOneByID(id string) (entities.JoinOrganizationRequest, error)
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

func (service *joinOrganizationRequestsImpl) Create(userID, organizationID string) (entities.JoinOrganizationRequest, error) {
	data := entities.JoinOrganizationRequest{
		UserID:         userID,
		OrganizationID: organizationID,
	}
	return service.repository.Create(data)
}

func (service *joinOrganizationRequestsImpl) FindManyByOrganizationID(organizationID string, offset, limit int64) (int64, []entities.JoinOrganizationRequest, error) {
	return service.repository.FindManyByOrganizationID(organizationID, offset, limit)
}

func (service *joinOrganizationRequestsImpl) FindManyByUserID(userID string, offset, limit int64) (int64, []entities.JoinOrganizationRequest, error) {
	return service.repository.FindManyByUserID(userID, offset, limit)
}

func (service *joinOrganizationRequestsImpl) FindOneByID(id string) (entities.JoinOrganizationRequest, error) {
	return service.repository.FindOneByID(id)
}

func (service *joinOrganizationRequestsImpl) Accept(id string, auditor entities.User) error {
	request, err := service.repository.FindOneByID(id)

	if err != nil {
		return err
	}

	user, err := service.usersService.FindOneByID(request.UserID)

	if err != nil {
		return err
	}

	request.AcceptedAt = time.Now()
	request.Status = utils.AcceptedStatus
	request.AuditorID = auditor.ID

	if err = service.repository.UpdateOneByID(request.ID, request); err != nil {
		return err
	}

	user.OrganizationID = request.OrganizationID
	user.PlatformRole = utils.OrgMemberPlatformRole

	if err = service.usersService.UpdateOneByID(user.ID, user); err != nil {
		return err
	}

	return nil
}

func (service *joinOrganizationRequestsImpl) Reject(id string, auditor entities.User, data entities.JoinOrganizationRequest) error {
	request, err := service.repository.FindOneByID(id)

	if err != nil {
		return err
	}

	request.AuditorID = auditor.ID
	request.Status = utils.RejectedStatus
	request.RejectedAt = time.Now()
	request.RejectReason = data.RejectReason

	if err = service.repository.UpdateOneByID(request.ID, request); err != nil {
		return err
	}

	return err
}
