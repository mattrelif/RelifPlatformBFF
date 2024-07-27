package services

import (
	"relif/bff/entities"
	"relif/bff/repositories"
	"time"
)

type JoinOrganizationRequests interface {
	Create(userID string, invite entities.JoinOrganizationRequest) (string, error)
	FindManyByOrganizationId(organizationId string, offset, limit int64) (int64, []entities.JoinOrganizationRequest, error)
	Accept(id string) error
	Reject(id string) error
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

func (service *joinOrganizationRequestsImpl) Create(userID string, request entities.JoinOrganizationRequest) (string, error) {
	request.UserID = userID
	request.CreatedAt = time.Now()

	return service.repository.Create(request)
}

func (service *joinOrganizationRequestsImpl) FindManyByOrganizationId(organizationId string, offset, limit int64) (int64, []entities.JoinOrganizationRequest, error) {
	return service.repository.FindManyByOrganizationId(organizationId, offset, limit)
}

func (service *joinOrganizationRequestsImpl) Accept(id string) error {
	request, err := service.repository.FindOneAndDeleteById(id)

	if err != nil {
		return err
	}

	data := entities.User{
		OrganizationID: request.OrganizationID,
	}
	if err = service.usersService.UpdateOneById(request.UserID, data); err != nil {
		return err
	}

	return nil
}

func (service *joinOrganizationRequestsImpl) Reject(id string) error {
	return service.repository.DeleteOneById(id)
}
